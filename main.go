package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"telegram-processor/cmd"
	"telegram-processor/internal/logger"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	logger.InitDefaulLog()

	// run app
	err := cmd.Execute(ctx)
	if errors.Is(err, context.Canceled) {
		slog.Warn("Command was interrupted", "error", err)
		os.Exit(0)
	}
	if err != nil {
		slog.Error("Failed to execute command", "error", err)
		os.Exit(1)
	}
}
