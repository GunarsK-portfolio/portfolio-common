package health

import (
	"context"
	"testing"
)

func TestNewPgxChecker(t *testing.T) {
	checker := NewPgxChecker(nil)

	if checker == nil {
		t.Fatal("expected checker to not be nil")
	}
}

func TestPgxChecker_Name(t *testing.T) {
	checker := NewPgxChecker(nil)

	if checker.Name() != "postgres" {
		t.Errorf("expected name 'postgres', got %s", checker.Name())
	}
}

func TestPgxChecker_Check_NilPool(t *testing.T) {
	checker := NewPgxChecker(nil)

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil pool, got %s", result.Status)
	}
	if result.Error != "pool is nil" {
		t.Errorf("expected 'pool is nil' error, got %s", result.Error)
	}
}
