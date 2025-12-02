package health

import (
	"context"
	"testing"
)

func TestNewPostgresChecker(t *testing.T) {
	checker := NewPostgresChecker(nil)

	if checker == nil {
		t.Fatal("expected checker to not be nil")
	}
}

func TestPostgresChecker_Name(t *testing.T) {
	checker := NewPostgresChecker(nil)

	if checker.Name() != "postgres" {
		t.Errorf("expected name 'postgres', got %s", checker.Name())
	}
}

func TestPostgresChecker_Check_NilDB(t *testing.T) {
	checker := NewPostgresChecker(nil)

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil db, got %s", result.Status)
	}
	if result.Error != "database is nil" {
		t.Errorf("expected 'database is nil' error, got %s", result.Error)
	}
}
