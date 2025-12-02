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
	// When db.DB() is called on nil, it will panic, so we test with
	// a properly initialized checker but this is more of an integration test scenario.
	// For unit testing, we verify the interface and name.
	checker := NewPostgresChecker(nil)

	// This will panic because we can't call methods on nil *gorm.DB
	// In real usage, the db is never nil
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when checking nil db")
		}
	}()

	_ = checker.Check(context.Background())
}
