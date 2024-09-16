package api

import (
	"context"
	"fmt"
	"log/slog"
	"telegram-processor/internal/models"
	pb "telegram-processor/pkg/api/proto"
)

const SEARCH_LIMIT = 10

func (s *processorApiServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	traceId, err := GetTraceIdFromCtx(ctx)
	if err != nil {
		slog.Warn("[Search]", "error", err)
	}
	slog.Debug("Search request", "query", in.Query, "trace-id", traceId)

	var messages models.Messages
	messages, err = s.processor.GetClosest(ctx, in.Query, SEARCH_LIMIT)
	if err != nil {
		err = fmt.Errorf("s.processor.GetClosest -> %w", err)
		slog.Error("Failed to search messages", "error", err)
		return nil, ErrSomethingWentWrong
	}

	if len(messages) == 0 {
		return nil, ErrNotFound
	}

	return &pb.SearchResponse{Messages: messages.ToPbMessageSearchedSimple()}, nil
}
