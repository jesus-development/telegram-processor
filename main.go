package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"telegram-processor/cmd"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

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
