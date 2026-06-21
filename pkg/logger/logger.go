package logger

import (
	"log/slog"
	"os"
	"strings"
)

var logLevel = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// New - creates a new instance of the slog logger,
// and sets it to be the default handler logger.
//   - level - The level that we want to logs to log at, e.g. debug, info, warn, error
func New(level string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel[strings.ToLower(level)],
	}))
	slog.SetDefault(logger)
}
