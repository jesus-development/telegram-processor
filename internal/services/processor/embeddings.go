package processor

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"telegram-processor/internal/models"
)

func (p *messageProcessor) CalculateAndSaveEmbeddings(ctx context.Context) error {
	messages, err := p.MessagesRepository.GetMessagesWithoutVectors(ctx)
	if err != nil {
		return fmt.Errorf("p.MessagesRepository.GetMessagesWithoutVectors -> %w", err)
	}

	embeddings, err := p.EmbeddingService.GetMessageEmbeddings(ctx, messages)
	if err != nil {
		return fmt.Errorf("p.EmbeddingService.GetMessageEmbeddings -> %w", err)
	}

	if err = p.MessagesRepository.InsertEmbeddings(ctx, embeddings); err != nil {
		return fmt.Errorf("p.MessagesRepository.InsertEmbeddings -> %w", err)
	}

	return nil
}

func (p *messageProcessor) GetEmbeddingPrice(ctx context.Context, tarif models.EmbeddingTarif) (decimal.Decimal, error) {
	messages, err := p.MessagesRepository.GetMessagesWithoutVectors(ctx)
	if err != nil {
		return decimal.Zero, fmt.Errorf("p.MessagesRepository.GetMessagesWithoutVectors -> %w", err)
	}

	price, err := p.EmbeddingService.GetEmbeddingsPrice(ctx, messages, tarif)
	if err != nil {
		return decimal.Zero, fmt.Errorf("p.EmbeddingService.GetEmbeddingsPrice -> %w", err)
	}

	return price, nil
}

func (p *messageProcessor) GetClosest(ctx context.Context, search string, limit int64) ([]*models.Message, error) {
	embedding, err := p.EmbeddingService.GetEmbedding(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("p.EmbeddingService.GetEmbedding -> %w", err)
	}

	messages, err := p.MessagesRepository.GetClosest(ctx, embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("p.MessagesRepository.GetClosest -> %w", err)
	}

	return messages, nil
}
