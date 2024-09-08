package json

import "strings"

type (
	ChatJson struct {
		ID       int64          `json:"id"`
		Type     string         `json:"type"`
		Messages []*MessageJson `json:"messages"`
	}

	MessageJson struct {
		TelegramID   int64             `json:"id"`
		TextEntities []*TextEntityJson `json:"text_entities"`
		Username     string            `json:"from"`
		UserID       string            `json:"from_id"`
		DateUnix     string            `json:"date_unixtime"`
	}

	TextEntityJson struct {
		Text string `json:"text"`
	}
)

func (m *MessageJson) String() string {
	var sb strings.Builder
	for _, te := range m.TextEntities {
		sb.WriteString(te.Text)
	}
	return sb.String()
}
