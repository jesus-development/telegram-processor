package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"net"
	"net/http"
	"telegram-processor/internal/services/processor"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "telegram-processor/pkg/api/proto"
)

const HTTP_SHUTDOWN_TIMEOUT = 10 * time.Second

type processorApiServer struct {
	pb.UnimplementedTelegramProcessorServiceServer
	grpcServer *grpc.Server
	httpServer *http.Server

	processor processor.MessageProcessor
}

func NewServer(processor processor.MessageProcessor) *processorApiServer {
	return &processorApiServer{processor: processor}
}

func (srv *processorApiServer) ListenGRPC() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("net.Listen -> %w", err)
	}

	// todo metadata, middlewares
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"server-metadata": []string{"value"}})

	srv.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(ctx, req)
	}))

	reflection.Register(srv.grpcServer)
	pb.RegisterTelegramProcessorServiceServer(srv.grpcServer, srv)

	slog.Info("Serving gRPC on 0.0.0.0:50051")
	err = srv.grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("srv.grpcServer.Serve -> %w", err)
	}

	return nil
}

// ListenHTTPGateway Must be called after ListenGRPC
func (srv *processorApiServer) ListenHTTPGateway() error {
	conn, err := grpc.NewClient(
		"0.0.0.0:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("grpc.NewClient -> %w", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterTelegramProcessorServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		return fmt.Errorf("pb.RegisterTelegramProcessorServiceHandler -> %w", err)
	}

	srv.httpServer = &http.Server{
		Addr:    ":50052",
		Handler: gwmux,
	}

	slog.Info("Serving http-gateway on http://0.0.0.0:50052")
	err = srv.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("srv.httpServer.ListenAndServe -> %w", err)
	}

	return nil
}

func (srv *processorApiServer) Shutdown() {
	if srv.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), HTTP_SHUTDOWN_TIMEOUT)
		err := srv.httpServer.Shutdown(ctx)
		if err != nil {
			slog.Error("Failed to shutdown http server", "error", err)
		}
	}
	if srv.grpcServer != nil {
		srv.grpcServer.GracefulStop()
	}
}
