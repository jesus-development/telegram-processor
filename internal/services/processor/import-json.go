package processor

import (
	"context"
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"telegram-processor/pkg/models/json"
)

func (p *messageProcessor) ImportJson(ctx context.Context, r io.Reader) error {
	var chat = &json.ChatJson{}
	err := easyjson.UnmarshalFromReader(r, chat)
	if err != nil {
		return fmt.Errorf("easyjson.UnmarshalFromReader -> %w", err)
	}

	if len(chat.Messages) == 0 {
		return ErrEmptyChat
	}

	if err = p.MessagesRepository.ImportChatFromJson(ctx, chat); err != nil {
		return fmt.Errorf("p.MessagesRepository.ImportChatFromJson -> %w", err)
	}

	return nil
}
