package logger

import (
	"fmt"
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"strings"
	"time"
)

const DEFAULT_LOG_LEVEL = slog.LevelDebug

var (
	lvl = new(slog.LevelVar)

	levelMap = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)

func InitDefaulLog() {
	lvl.Set(DEFAULT_LOG_LEVEL)
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      lvl,
			TimeFormat: time.DateTime,
		}),
	))
}

func SetLogLevel(level string) error {
	l, ok := levelMap[strings.ToLower(level)]
	if !ok {
		return fmt.Errorf("SetLogLevel -> %w", ErrInvalidLogLevel)
	}
	lvl.Set(l)

	return nil
}
