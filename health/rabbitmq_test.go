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

func TestNewRabbitMQCheckerWithProvider_NilProvider(t *testing.T) {
	checker := NewRabbitMQCheckerWithProvider(nil)

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil provider, got %s", result.Status)
	}
}

func TestNewRabbitMQCheckerWithProvider_Name(t *testing.T) {
	checker := NewRabbitMQCheckerWithProvider(nil)

	if checker.Name() != "rabbitmq" {
		t.Errorf("expected name 'rabbitmq', got %s", checker.Name())
	}
}

func TestQueueDepthChecker_Name(t *testing.T) {
	checker := NewQueueDepthChecker(nil, "contact_messages_dlq", 0)

	want := "queue:contact_messages_dlq"
	if checker.Name() != want {
		t.Errorf("expected name %q, got %s", want, checker.Name())
	}
}

func TestQueueDepthChecker_Check_NilConnection(t *testing.T) {
	checker := NewQueueDepthChecker(nil, "contact_messages_dlq", 0)

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil connection, got %s", result.Status)
	}
	if result.Error != "connection unavailable" {
		t.Errorf("expected 'connection unavailable' error, got %s", result.Error)
	}
}
