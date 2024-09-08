package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type (
	MessageEmbedding struct {
		ID        int64
		MessageID int64
		Embedding []float32
		CreatedAt time.Time
	}

	EmbeddingTarif struct {
		Price   decimal.Decimal
		Divider int64
	}
)
