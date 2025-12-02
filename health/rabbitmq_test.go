package health

import (
	"context"
	"testing"
)

func TestNewRabbitMQChecker(t *testing.T) {
	checker := NewRabbitMQChecker(nil)

	if checker == nil {
		t.Fatal("expected checker to not be nil")
	}
}

func TestRabbitMQChecker_Name(t *testing.T) {
	checker := NewRabbitMQChecker(nil)

	if checker.Name() != "rabbitmq" {
		t.Errorf("expected name 'rabbitmq', got %s", checker.Name())
	}
}

func TestRabbitMQChecker_Check_NilConnection(t *testing.T) {
	checker := NewRabbitMQChecker(nil)

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil connection, got %s", result.Status)
	}
	if result.Error != "connection is nil" {
		t.Errorf("expected 'connection is nil' error, got %s", result.Error)
	}
}
