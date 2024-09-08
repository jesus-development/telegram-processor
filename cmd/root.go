package cmd

import (
	"context"
	"github.com/spf13/cobra"
)

const ERR_PREFIX = "[cmd run]"

var rootCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "telegram-processor",
	Short:         "A CLI tool for processing Telegram messages.",
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
