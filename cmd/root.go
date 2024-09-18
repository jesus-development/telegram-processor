package cmd

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"telegram-processor/internal/config"
	"telegram-processor/internal/logger"
)

const ERR_PREFIX = "[cmd run]"

var (
	rootCmd = &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
		Use:           "telegram-processor",
		Short:         "A CLI tool for processing Telegram messages.",
	}

	appConfig config.Config

	// Flags
	cfgFile string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		config.DEFAULT_CONFIG_FILE,
		fmt.Sprintf("config file (default is %s)", config.DEFAULT_CONFIG_FILE))

	// Set up Cobra to initialize Viper and log level before executing the command
	cobra.OnInitialize(initConfig, setLogLevel)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		slog.Error("No config")
		os.Exit(1)
	}

	err := godotenv.Load(config.ENV_FILE)
	if err != nil {
		slog.Info("Error loading .env file")
	}

	// Secrets
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("openai.apikey", "OPENAI_API_KEY")

	if err = viper.ReadInConfig(); err != nil {
		slog.Error("Can't read config:", err)
		os.Exit(1)
	}

	if err = viper.Unmarshal(&appConfig); err != nil {
		slog.Error("Can't unmarshal config:", err)
		os.Exit(1)
	}
}

func setLogLevel() {
	if appConfig.LogLevel == "" {
		return
	}
	if err := logger.SetLogLevel(appConfig.LogLevel); err != nil {
		slog.Error("Can't set log level:", err)
	}
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
