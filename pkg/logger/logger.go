package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type contextKey string

// RequestIDKey is the context key for request ID.
const RequestIDKey contextKey = "request_id"

var defaultLogger *slog.Logger

// WithRequestID returns a new context with the request ID stored.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func Init(env string) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Ensure timestamp is always in RFC3339 format
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.String(slog.TimeKey, t.Format(time.RFC3339))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	if env == "production" {
		opts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	defaultLogger = slog.New(handler)
}

// Default returns the singleton logger instance.
func Default() *slog.Logger {
	if defaultLogger == nil {
		Init("development")
	}
	return defaultLogger
}

func withRequestID(ctx context.Context) *slog.Logger {
	log := Default()
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		return log.With(slog.String("request_id", requestID))
	}
	return log
}

// Info logs at INFO level with request_id from context.
func Info(ctx context.Context, msg string, args ...any) {
	withRequestID(ctx).Info(msg, args...)
}

// Error logs at ERROR level with request_id from context.
func Error(ctx context.Context, msg string, args ...any) {
	withRequestID(ctx).Error(msg, args...)
}

// Debug logs at DEBUG level with request_id from context.
func Debug(ctx context.Context, msg string, args ...any) {
	withRequestID(ctx).Debug(msg, args...)
}

// Warn logs at WARN level with request_id from context.
func Warn(ctx context.Context, msg string, args ...any) {
	withRequestID(ctx).Warn(msg, args...)
}
