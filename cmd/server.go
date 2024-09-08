//go:build api_server

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"runtime/debug"
	"sync"
	"telegram-processor/internal/api"
	"telegram-processor/internal/config"
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

		cfg := config.LoadConfig()

		db, err := database.NewDatabase(&cfg.DB)
		if err != nil {
			return fmt.Errorf("%s %s -> %w", ERR_PREFIX, cmd.Use, err)
		}
		messageRepo := messages.NewPGMessagesRepository(db)

		openaiService := openai.NewOpenAIService(&cfg.Openai)

		proc := processor.NewMessageProcessor(processor.WithMessagesRepository(messageRepo), processor.WithEmbeddingService(openaiService))

		apiServer := api.NewServer(proc)

		debug.Stack()
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
