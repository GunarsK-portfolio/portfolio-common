package health

import (
	"context"
	"testing"
)

func TestNewMinIOChecker(t *testing.T) {
	checker := NewMinIOChecker(nil, "test-bucket")

	if checker == nil {
		t.Fatal("expected checker to not be nil")
	}
}

func TestMinIOChecker_Name(t *testing.T) {
	checker := NewMinIOChecker(nil, "test-bucket")

	if checker.Name() != "minio" {
		t.Errorf("expected name 'minio', got %s", checker.Name())
	}
}

func TestMinIOChecker_Check_NilClient(t *testing.T) {
	checker := NewMinIOChecker(nil, "test-bucket")

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil client, got %s", result.Status)
	}
	if result.Error != "client is nil" {
		t.Errorf("expected 'client is nil' error, got %s", result.Error)
	}
}

func TestMinIOChecker_Check_NilClientNoBucket(t *testing.T) {
	checker := NewMinIOChecker(nil, "")

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil client, got %s", result.Status)
	}
	if result.Error != "client is nil" {
		t.Errorf("expected 'client is nil' error, got %s", result.Error)
	}
}
