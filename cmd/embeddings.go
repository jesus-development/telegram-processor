//go:build !api_server

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"telegram-processor/internal/config"
	database "telegram-processor/internal/db"
	"telegram-processor/internal/repository/messages"
	"telegram-processor/internal/services/external/openai"
	"telegram-processor/internal/services/processor"
)

var (
	embeddingsCmd = &cobra.Command{
		Use:   "embeddings",
		Short: "Embedding commands",
	}
	embeddingsCalcAndSaveCmd = &cobra.Command{
		Use:   "calc-and-save",
		Short: "Calculate and save embeddings",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			cfg := config.LoadConfig()

			db, err := database.NewDatabase(&cfg.DB)
			if err != nil {
				return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
			}

			messageRepo := messages.NewPGMessagesRepository(db)

			openaiService := openai.NewOpenAIService(&cfg.Openai)

			processor := processor.NewMessageProcessor(
				processor.WithMessagesRepository(messageRepo),
				processor.WithEmbeddingService(openaiService),
			)

			if err := processor.CalculateAndSaveEmbeddings(ctx); err != nil {
				return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
			}
			return nil
		},
	}

	embeddingsGetPriceCmd = &cobra.Command{
		Use:   "get-price",
		Short: "Get embedding price",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			cfg := config.LoadConfig()

			db, err := database.NewDatabase(&cfg.DB)
			if err != nil {
				return fmt.Errorf("%s %s database.NewDatabase -> %w", ERR_PREFIX, cmd.Use, err)
			}
			messageRepo := messages.NewPGMessagesRepository(db)

			openaiService := openai.NewOpenAIService(&cfg.Openai)

			processor := processor.NewMessageProcessor(
				processor.WithMessagesRepository(messageRepo),
				processor.WithEmbeddingService(openaiService),
			)

			price, err := processor.GetEmbeddingPrice(ctx, openai.DefaultTarif)
			if err != nil {
				return fmt.Errorf("%s %s processor.GetEmbeddingPrice -> %w", ERR_PREFIX, cmd.Use, err)
			}

			fmt.Printf("Price is $%s", price.String())

			return nil
		},
	}
)

func init() {
	embeddingsCmd.AddCommand(embeddingsCalcAndSaveCmd)
	embeddingsCmd.AddCommand(embeddingsGetPriceCmd)
	rootCmd.AddCommand(embeddingsCmd)
}
