//go:build demo

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strings"
	database "telegram-processor/internal/db"
	"telegram-processor/internal/repository/messages"
	"telegram-processor/internal/scenarios/demo"
	"telegram-processor/internal/services/external/openai"
	"telegram-processor/internal/services/processor"
	"telegram-processor/pkg/cli/prompt"
)

var (
	demoCmd = &cobra.Command{
		Use:   "demo",
		Short: "Demo commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			prompter := prompt.NewStdPrompter(os.Stdin, os.Stdout)

			// OpenAI API token check
			if appConfig.Openai.ApiKey == "" && !appConfig.Openai.IsFake {
				yes, err := prompter.YesNoPrompt("OpenAI API key is not set. Do you want to set it? Otherwise embeddings will be random", false)
				if err != nil {
					slog.Error("Can't get answer", "error", err)
				}

				if yes {
					token, err := prompter.StringPrompt("OpenAI API key:")
					if err != nil {
						slog.Error("Can't get answer", "error", err)
					} else {
						appConfig.Openai.ApiKey = strings.Trim(token, "\n")
					}
				}

				if appConfig.Openai.ApiKey == "" {
					slog.Info("Openai API key is not set. Embeddings and search result will be random.")
					appConfig.Openai.IsFake = true
				}
			}

			// Initialization
			var openaiService processor.EmbeddingService
			if appConfig.Openai.IsFake {
				openaiService = openai.NewFakeOpenAIService(&appConfig.Openai)
			} else {
				openaiService = openai.NewOpenAIService(&appConfig.Openai)
			}

			db, err := database.NewDatabase(&appConfig.DB)
			if err != nil {
				return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
			}

			messageRepo := messages.NewPGMessagesRepository(db)

			processor := processor.NewMessageProcessor(
				processor.WithMessagesRepository(messageRepo),
				processor.WithEmbeddingService(openaiService),
			)

			importDbFlag, err := cmd.Flags().GetBool("import-db")
			if err != nil {
				slog.Error("Can't read flag import-db", "error", err)
			}

			// Run demo
			err = demo.NewDemoScenario(processor, prompter, importDbFlag).Run(ctx)
			if err != nil {
				return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
			}

			return nil
		},
	}
)

func init() {
	demoCmd.Flags().Bool("import-db", false, "import-db")
	rootCmd.AddCommand(demoCmd)
}
