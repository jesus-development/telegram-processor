package processor

import (
	"context"
	"github.com/shopspring/decimal"
	"io"
	"telegram-processor/internal/models"
	"telegram-processor/pkg/models/json"
)

type (
	MessagesRepository interface {
		GetMessagesWithoutVectors(ctx context.Context) ([]*models.Message, error)
		ImportChatFromJson(ctx context.Context, chat *json.ChatJson) error
		InsertEmbeddings(ctx context.Context, embeddings []*models.MessageEmbedding) error
		GetClosest(ctx context.Context, search []float32, limit int64) ([]*models.Message, error)
		GetCount(ctx context.Context) (int64, error)
	}

	EmbeddingService interface {
		GetMessageEmbeddings(ctx context.Context, messages []*models.Message) ([]*models.MessageEmbedding, error)
		GetEmbeddingsPrice(ctx context.Context, messages []*models.Message, tarif models.EmbeddingTarif) (decimal.Decimal, error)
		GetEmbedding(ctx context.Context, text string) ([]float32, error)
	}

	MessageProcessor interface {
		// import-json
		ImportJson(ctx context.Context, reader io.Reader) error

		// embeddings
		CalculateAndSaveEmbeddings(ctx context.Context) error
		GetEmbeddingPrice(ctx context.Context, tarif models.EmbeddingTarif) (decimal.Decimal, error)
		GetClosest(ctx context.Context, search string, limit int64) ([]*models.Message, error)

		// messages
		GetCount(ctx context.Context) (int64, error)
	}

	messageProcessor struct {
		MessagesRepository MessagesRepository
		EmbeddingService   EmbeddingService
	}
)

func NewMessageProcessor(opts ...ProcessorOption) (mp *messageProcessor) {
	mp = &messageProcessor{}

	for _, opt := range opts {
		opt(mp)
	}
	return mp
}
