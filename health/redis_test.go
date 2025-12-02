package health

import (
	"context"
	"testing"
)

func TestNewRedisChecker(t *testing.T) {
	checker := NewRedisChecker(nil)

	if checker == nil {
		t.Fatal("expected checker to not be nil")
	}
}

func TestRedisChecker_Name(t *testing.T) {
	checker := NewRedisChecker(nil)

	if checker.Name() != "redis" {
		t.Errorf("expected name 'redis', got %s", checker.Name())
	}
}

func TestRedisChecker_Check_NilClient(t *testing.T) {
	checker := NewRedisChecker(nil)

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil client, got %s", result.Status)
	}
	if result.Error != "client is nil" {
		t.Errorf("expected 'client is nil' error, got %s", result.Error)
	}
}
