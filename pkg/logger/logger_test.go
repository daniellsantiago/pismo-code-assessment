package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	t.Run("logs with request_id when present in context", func(t *testing.T) {
		var buf bytes.Buffer
		defaultLogger = slog.New(slog.NewTextHandler(&buf, nil))

		ctx := context.WithValue(context.Background(), RequestIDKey, "test-request-id-123")

		Info(ctx, "test message")

		assert.Contains(t, buf.String(), "request_id=test-request-id-123")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("logs without request_id when not in context", func(t *testing.T) {
		var buf bytes.Buffer
		defaultLogger = slog.New(slog.NewTextHandler(&buf, nil))

		Info(context.Background(), "test message")

		assert.NotContains(t, buf.String(), "request_id")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("includes timestamp in logs", func(t *testing.T) {
		var buf bytes.Buffer
		defaultLogger = slog.New(slog.NewTextHandler(&buf, nil))

		Info(context.Background(), "test message")

		assert.Contains(t, buf.String(), "time=")
	})
}

func TestError(t *testing.T) {
	t.Run("logs error with request_id", func(t *testing.T) {
		var buf bytes.Buffer
		defaultLogger = slog.New(slog.NewTextHandler(&buf, nil))

		ctx := context.WithValue(context.Background(), RequestIDKey, "error-request-id")

		Error(ctx, "error occurred", slog.String("details", "something bad"))

		assert.Contains(t, buf.String(), "request_id=error-request-id")
		assert.Contains(t, buf.String(), "error occurred")
		assert.Contains(t, buf.String(), "details=")
	})
}

func TestInit(t *testing.T) {
	t.Run("creates JSON logger for production", func(t *testing.T) {
		Init("production")
		assert.NotNil(t, defaultLogger)
	})

	t.Run("creates text logger for development", func(t *testing.T) {
		Init("development")
		assert.NotNil(t, defaultLogger)
	})
}

func TestDefault(t *testing.T) {
	t.Run("returns default logger", func(t *testing.T) {
		defaultLogger = nil
		log := Default()
		assert.NotNil(t, log)
	})
}
