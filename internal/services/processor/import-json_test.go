package processor

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"io"
	"os"
	"telegram-processor/internal/services/processor/mock"
	"testing"
)

func TestImportJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	messageRepoMock := mock.NewMockMessagesRepository(ctrl)

	processor := NewMessageProcessor(WithMessagesRepository(messageRepoMock))

	emptyReader, err := os.Open("test-data/import-json_empty.json")
	if err != nil {
		t.Fatal(err)
	}
	defer emptyReader.Close()

	goodReader, err := os.Open("test-data/import-json_good.json")
	if err != nil {
		t.Fatal(err)
	}
	defer goodReader.Close()

	testCases := []struct {
		title    string
		mockFunc func()
		input    io.Reader
		err      error
	}{
		{
			title:    "Empty chat",
			mockFunc: func() {},
			input:    emptyReader,
			err:      ErrEmptyChat,
		},
		{
			title: "Success",
			mockFunc: func() {
				messageRepoMock.EXPECT().ImportChatFromJson(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: goodReader,
			err:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			tc.mockFunc()
			err := processor.ImportJson(context.Background(), tc.input)

			if !errors.Is(err, tc.err) {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
		})
	}
}
