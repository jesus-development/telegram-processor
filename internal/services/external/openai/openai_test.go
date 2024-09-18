package openai

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"telegram-processor/internal/models"
	"testing"
)

func TestGetEmbeddingsPrice(t *testing.T) {
	t.Parallel()
	// 10 short words = 10 tokens
	messages := []*models.Message{
		{Text: "Hello"},
		{Text: "World"},
		{Text: "Hello World"},
		{Text: "Hello World"},
		{Text: "Hello World Hello World"},
	}

	testCases := []struct {
		title    string
		tarif    models.EmbeddingTarif
		messages []*models.Message
		err      error
		expected decimal.Decimal
	}{
		{
			title:    "Success. 10 tokens with tarif (100$ per 10 tokens)",
			tarif:    models.EmbeddingTarif{Price: decimal.NewFromFloat(100), Divider: 10},
			messages: messages,
			err:      nil,
			expected: decimal.NewFromInt(100),
		},
		{
			title:    "Bad tarif",
			tarif:    models.EmbeddingTarif{Price: decimal.NewFromFloat(1.23), Divider: 0},
			messages: []*models.Message{{Text: "test"}},
			err:      ErrInvalidTarif,
			expected: decimal.Zero,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			s := &openAIService{}
			actual, err := s.GetEmbeddingsPrice(context.Background(), tc.messages, tc.tarif)
			if !errors.Is(err, tc.err) {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if actual.Cmp(tc.expected) != 0 {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
