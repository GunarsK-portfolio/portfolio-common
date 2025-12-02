package queue

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

// =============================================================================
// GetRetryCount Tests
// =============================================================================

func TestGetRetryCount(t *testing.T) {
	tests := []struct {
		name     string
		headers  amqp.Table
		expected int
	}{
		{
			name:     "nil headers returns 0",
			headers:  nil,
			expected: 0,
		},
		{
			name:     "empty headers returns 0",
			headers:  amqp.Table{},
			expected: 0,
		},
		{
			name:     "missing header returns 0",
			headers:  amqp.Table{"other-header": "value"},
			expected: 0,
		},
		{
			name:     "header with int32 value",
			headers:  amqp.Table{RetryCountHeader: int32(3)},
			expected: 3,
		},
		{
			name:     "header with int64 value",
			headers:  amqp.Table{RetryCountHeader: int64(5)},
			expected: 5,
		},
		{
			name:     "header with int value",
			headers:  amqp.Table{RetryCountHeader: int(2)},
			expected: 2,
		},
		{
			name:     "header with zero value",
			headers:  amqp.Table{RetryCountHeader: int32(0)},
			expected: 0,
		},
		{
			name:     "header with string value returns 0",
			headers:  amqp.Table{RetryCountHeader: "invalid"},
			expected: 0,
		},
		{
			name:     "header with float value returns 0",
			headers:  amqp.Table{RetryCountHeader: 3.14},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delivery := amqp.Delivery{Headers: tt.headers}
			result := GetRetryCount(delivery)

			if result != tt.expected {
				t.Errorf("GetRetryCount() = %d, want %d", result, tt.expected)
			}
		})
	}
}

// =============================================================================
// RetryCountHeader Tests
// =============================================================================

func TestRetryCountHeader_Value(t *testing.T) {
	want := "x-retry-count"
	if RetryCountHeader != want {
		t.Errorf("RetryCountHeader = %q, want %q", RetryCountHeader, want)
	}
}

// =============================================================================
// Error Definitions Tests
// =============================================================================

func TestConsumerErrorDefinitions(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantMsg string
	}{
		{
			name:    "ErrConsumerClosed",
			err:     ErrConsumerClosed,
			wantMsg: "consumer is closed",
		},
		{
			name:    "ErrConsumeSetupFailed",
			err:     ErrConsumeSetupFailed,
			wantMsg: "failed to setup consumer",
		},
		{
			name:    "ErrNilPublisher",
			err:     ErrNilPublisher,
			wantMsg: "publisher is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.wantMsg {
				t.Errorf("%s.Error() = %q, want %q", tt.name, tt.err.Error(), tt.wantMsg)
			}
		})
	}
}

// =============================================================================
// Close Tests
// =============================================================================

func TestConsumerClose_Idempotent(t *testing.T) {
	consumer := &RabbitMQConsumer{
		closed: false,
	}

	// First close
	err := consumer.Close()
	if err != nil {
		t.Errorf("First Close() error = %v", err)
	}

	// Second close should also succeed (idempotent)
	err = consumer.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v, want nil (idempotent)", err)
	}

	if !consumer.closed {
		t.Error("Consumer should be marked as closed")
	}
}

func TestConsumerClose_AlreadyClosed(t *testing.T) {
	consumer := &RabbitMQConsumer{
		closed: true,
	}

	err := consumer.Close()
	if err != nil {
		t.Errorf("Close() on already closed consumer error = %v, want nil", err)
	}
}

// =============================================================================
// Interface Compliance Tests
// =============================================================================

func TestConsumerInterfaceCompliance(t *testing.T) {
	var _ Consumer = (*RabbitMQConsumer)(nil)
}
