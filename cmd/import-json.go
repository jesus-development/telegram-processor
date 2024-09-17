//go:build tools

package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	database "telegram-processor/internal/db"
	"telegram-processor/internal/repository/messages"
	"telegram-processor/internal/services/processor"
)

var ErrNoFileProvided = errors.New("no file path provided")

var importJsonCmd = &cobra.Command{
	Use:   "import-json",
	Short: "Load and process Telegram messages",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		db, err := database.NewDatabase(&appConfig.DB)
		if err != nil {
			return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
		}
		messageRepo := messages.NewPGMessagesRepository(db)

		processor := processor.NewMessageProcessor(processor.WithMessagesRepository(messageRepo))

		// Parse JSON and save messages to repo
		if len(args) == 0 {
			return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, ErrNoFileProvided)
		}

		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
		}
		defer func() {
			if err = file.Close(); err != nil {
				slog.Error("Failed to close file", "error", err)
			}
		}()

		if err = processor.ImportJson(ctx, file); err != nil {
			return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importJsonCmd)
}
