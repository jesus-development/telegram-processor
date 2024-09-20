package openai

import (
	"context"
	"math/rand/v2"
	"telegram-processor/internal/config"
	"telegram-processor/internal/models"
)

type (
	fakeOpenAIService struct {
		openAIService
	}
)

func NewFakeOpenAIService(cfg *config.OpenaiConfig) *fakeOpenAIService {
	return &fakeOpenAIService{}
}

func (s *fakeOpenAIService) GetMessageEmbeddings(ctx context.Context, messages []*models.Message) ([]*models.MessageEmbedding, error) {
	messageEmbeddings := make([]*models.MessageEmbedding, 0, len(messages))
	for _, msg := range messages {
		messageEmbedding := &models.MessageEmbedding{
			Embedding: randomFloat32Slice(MAX_EMBEDDINGS_DIMENSIONS),
			MessageID: msg.ID,
		}
		messageEmbeddings = append(messageEmbeddings, messageEmbedding)
	}

	return messageEmbeddings, nil
}

func (s *fakeOpenAIService) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	return randomFloat32Slice(MAX_EMBEDDINGS_DIMENSIONS), nil
}

func randomFloat32Slice(size int) []float32 {
	slice := make([]float32, size)
	for i := range slice {
		slice[i] = rand.Float32()
	}
	return slice
}
