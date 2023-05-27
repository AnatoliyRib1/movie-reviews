package log

import (
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	"os"
)

var contextKey = struct{}{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(contextKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func SetupLogger(islocal bool, level string) (*slog.Logger, error) {
	l, err := newLevelFromString(level)
	if err != nil {
		return nil, err
	}
	opts := &slog.HandlerOptions{Level: l}

	var handler slog.Handler

	if islocal {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler), nil
}

func newLevelFromString(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	}
	return slog.Level(0), fmt.Errorf("unknown level :%q", level)
}
