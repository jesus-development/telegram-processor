package models

import (
	pbmodel "telegram-processor/pkg/models/proto"
	"time"
)

type (
	Message struct {
		ID         int64
		TelegramID int64
		Text       string
		Length     int64
		UserID     int64
		Username   string
		Date       time.Time
		CreatedAt  time.Time
		ChatID     int64
		Similarity float32
	}

	Messages []*Message
)

func (m Messages) ToPbMessageSearchedSimple() []*pbmodel.MessageSearchedSimple {
	res := make([]*pbmodel.MessageSearchedSimple, len(m))

	for i, msg := range m {
		res[i] = &pbmodel.MessageSearchedSimple{
			Text:       msg.Text,
			Similarity: msg.Similarity,
		}
	}

	return res
}
