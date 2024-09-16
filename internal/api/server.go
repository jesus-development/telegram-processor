package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"net"
	"net/http"
	"strings"
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

	// todo middlewares
	srv.grpcServer = grpc.NewServer()

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

	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcher),
	)

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

func headerMatcher(key string) (string, bool) {
	switch key {
	case "X-Trace-Id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func GetTraceIdFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata.FromIncomingContext -> %w", ErrGetTraceId)
	}
	values := md["x-trace-id"]
	if len(values) == 0 {
		return "", fmt.Errorf("len(values) == 0 -> %w", ErrGetTraceId)
	}

	return strings.Join(values, ""), nil
}
