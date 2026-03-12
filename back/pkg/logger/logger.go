package logger

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

func New(cfg Config) (*slog.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	options := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     level,
	}

	switch strings.ToLower(strings.TrimSpace(cfg.Format)) {
	case "text":
		return slog.New(slog.NewTextHandler(os.Stdout, options)), nil
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, options)), nil
	default:
		return nil, errors.New("unsupported log format")
	}
}

func parseLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, errors.New("unsupported log level")
	}
}
