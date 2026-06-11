package logger

import (
	"context"
	"log/slog"
	"testing"
)

func TestNewFromEnv_Defaults(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("LOG_FORMAT", "")
	t.Setenv("LOG_SOURCE", "")

	log := NewFromEnv("test-service")
	if !log.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("default level should enable info")
	}
	if log.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("default level should not enable debug")
	}
}

func TestNewFromEnv_ReadsLevel(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")

	log := NewFromEnv("test-service")
	if !log.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("LOG_LEVEL=debug should enable debug")
	}
}

func TestNewFromEnv_InvalidLogSourcePanics(t *testing.T) {
	t.Setenv("LOG_SOURCE", "yes")

	defer func() {
		if recover() == nil {
			t.Error("expected panic for invalid LOG_SOURCE value")
		}
	}()
	NewFromEnv("test-service")
}
