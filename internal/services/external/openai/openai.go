package openai

import (
	"context"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"github.com/shopspring/decimal"
	"telegram-processor/internal/config"
	"telegram-processor/internal/models"
)

type (
	openAIService struct {
		client *openai.Client
	}
)

const (
	// 2000 is max for index
	MAX_EMBEDDINGS_DIMENSIONS = 2000

	// $ per 1 million tokens
	TARIF_LARGE3_PRICE   = 0.13
	TARIF_LARGE3_DIVIDER = int64(1_000_000)
)

var DefaultTarif = models.EmbeddingTarif{Price: decimal.NewFromFloat(TARIF_LARGE3_PRICE), Divider: TARIF_LARGE3_DIVIDER}

func NewOpenAIService(cfg *config.OpenaiConfig) *openAIService {
	return &openAIService{
		client: openai.NewClient(cfg.ApiKey),
	}
}

func (s *openAIService) GetMessageEmbeddings(ctx context.Context, messages []*models.Message) ([]*models.MessageEmbedding, error) {
	if len(messages) == 0 {
		return make([]*models.MessageEmbedding, 0), nil
	}

	input := make([]string, len(messages))
	for i, message := range messages {
		input[i] = message.Text
	}
	req := openai.EmbeddingRequestStrings{
		Input:      input,
		Model:      openai.LargeEmbedding3,
		Dimensions: MAX_EMBEDDINGS_DIMENSIONS,
	}

	resp, err := s.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("s.client.CreateEmbeddings -> %w", err)
	}

	messageEmbeddings := make([]*models.MessageEmbedding, 0, len(resp.Data))
	for _, emb := range resp.Data {
		messageEmbedding := &models.MessageEmbedding{
			Embedding: emb.Embedding,
			MessageID: messages[emb.Index].ID,
		}
		messageEmbeddings = append(messageEmbeddings, messageEmbedding)
	}

	return messageEmbeddings, nil
}

func (s *openAIService) GetEmbeddingsPrice(ctx context.Context, messages []*models.Message, tarif models.EmbeddingTarif) (decimal.Decimal, error) {
	if tarif.Divider < 1 {
		return decimal.Zero, fmt.Errorf("%w: bad divider", ErrInvalidTarif)
	}

	tkm, err := tiktoken.EncodingForModel(string(openai.LargeEmbedding3))
	if err != nil {
		return decimal.Zero, fmt.Errorf("tiktoken.EncodingForModel -> %w", err)
	}

	tokensCount := 0
	for _, msg := range messages {
		t := tkm.Encode(msg.Text, nil, nil)
		tokensCount += len(t)
	}

	return tarif.Price.Mul(decimal.NewFromInt(int64(tokensCount))).Div(decimal.NewFromInt(tarif.Divider)), nil
}

func (s *openAIService) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input:      text,
		Model:      openai.LargeEmbedding3,
		Dimensions: MAX_EMBEDDINGS_DIMENSIONS,
	}

	resp, err := s.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("s.client.CreateEmbeddings -> %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, ErrNoEmbedding
	}

	return resp.Data[0].Embedding, nil
}
