//go:build api_server

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"sync"
	"telegram-processor/internal/api"
	database "telegram-processor/internal/db"
	"telegram-processor/internal/repository/messages"
	"telegram-processor/internal/services/external/openai"
	"telegram-processor/internal/services/processor"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start gRPC and http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		db, err := database.NewDatabase(&appConfig.DB)
		if err != nil {
			return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
		}
		messageRepo := messages.NewPGMessagesRepository(db)

		var openaiService processor.EmbeddingService
		// todo (appConfig.Openai.ApiKey == "") will be removed soon
		if appConfig.Openai.ApiKey == "" || appConfig.Openai.IsFake {
			slog.Warn("Openai API key is not set. Embeddings and search result will be random.")
			openaiService = openai.NewFakeOpenAIService(&appConfig.Openai)
		} else {
			openaiService = openai.NewOpenAIService(&appConfig.Openai)
		}

		proc := processor.NewMessageProcessor(processor.WithMessagesRepository(messageRepo), processor.WithEmbeddingService(openaiService))

		apiServer := api.NewServer(proc, &appConfig.Server)

		var (
			chErr = make(chan error, 2)
			wg    = &sync.WaitGroup{}
		)

		wg.Add(2)
		go func() {
			if err := apiServer.ListenGRPC(); err != nil {
				chErr <- fmt.Errorf("apiServer.ListenGRPC -> %w", err)
			}

			slog.Info("gRPC server stopped")
			wg.Done()
		}()

		go func() {
			if err := apiServer.ListenHTTPGateway(); err != nil {
				chErr <- fmt.Errorf("apiServer.ListenHTTPGateway -> %w", err)
			}

			slog.Info("http-gateway server stopped")
			wg.Done()
		}()

		select {
		case err := <-chErr:
			apiServer.Shutdown()
			wg.Wait()
			return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
		case <-ctx.Done():
			apiServer.Shutdown()
			wg.Wait()
			return ctx.Err()
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
