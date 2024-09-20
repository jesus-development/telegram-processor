package processor

import (
	"context"
)

func (p *messageProcessor) GetCount(ctx context.Context) (int64, error) {
	return p.MessagesRepository.GetCount(ctx)
}
