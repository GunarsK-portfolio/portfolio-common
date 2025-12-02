package queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/GunarsK-portfolio/portfolio-common/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RetryCountHeader is the AMQP header key for tracking retry attempts
const RetryCountHeader = "x-retry-count"

// Consumer errors
var (
	ErrConsumerClosed     = errors.New("consumer is closed")
	ErrConsumeSetupFailed = errors.New("failed to setup consumer")
	ErrNilPublisher       = errors.New("publisher is required")
)

// MessageHandler processes a single message delivery.
// Return nil to ACK the message, return error to trigger retry logic.
type MessageHandler func(ctx context.Context, delivery amqp.Delivery) error

// Consumer defines the interface for message queue consuming
type Consumer interface {
	// Consume starts consuming messages and blocks until context is cancelled
	Consume(ctx context.Context, handler MessageHandler) error
	// Close stops consuming and closes connections
	Close() error
}

// RabbitMQConsumer implements Consumer for RabbitMQ
type RabbitMQConsumer struct {
	mu        sync.Mutex
	closed    bool
	conn      *amqp.Connection
	channel   *amqp.Channel
	publisher *RabbitMQPublisher
	config    config.RabbitMQConfig
	logger    *slog.Logger
}

// NewRabbitMQConsumer creates a new consumer that shares queue infrastructure with the publisher.
// The publisher must be created first as it declares all queues.
func NewRabbitMQConsumer(
	cfg config.RabbitMQConfig,
	publisher *RabbitMQPublisher,
	logger *slog.Logger,
) (*RabbitMQConsumer, error) {
	if publisher == nil {
		return nil, ErrNilPublisher
	}

	conn, err := amqp.Dial(cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("%w: %v", ErrChannelFailed, err)
	}

	cleanup := func() {
		_ = ch.Close()
		_ = conn.Close()
	}

	// Set QoS for fair dispatch
	if err := ch.Qos(cfg.PrefetchCount, 0, false); err != nil {
		cleanup()
		return nil, fmt.Errorf("%w: set qos: %v", ErrConsumeSetupFailed, err)
	}

	return &RabbitMQConsumer{
		conn:      conn,
		channel:   ch,
		publisher: publisher,
		config:    cfg,
		logger:    logger,
	}, nil
}

// Consume starts consuming messages from the queue.
// Blocks until context is cancelled or an error occurs.
// The handler is called for each message; return nil to ACK, error to handle retry.
func (c *RabbitMQConsumer) Consume(ctx context.Context, handler MessageHandler) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return ErrConsumerClosed
	}
	c.mu.Unlock()

	deliveries, err := c.channel.Consume(
		c.config.Queue,
		c.config.ConsumerTag,
		false, // auto-ack disabled for manual control
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("%w: consume: %v", ErrConsumeSetupFailed, err)
	}

	c.logger.Info("Consumer started", "queue", c.config.Queue, "tag", c.config.ConsumerTag)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Consumer stopping", "reason", ctx.Err())
			return ctx.Err()

		case delivery, ok := <-deliveries:
			if !ok {
				c.logger.Warn("Delivery channel closed")
				return errors.New("delivery channel closed")
			}

			c.processDelivery(ctx, delivery, handler)
		}
	}
}

// processDelivery handles a single message with retry logic
func (c *RabbitMQConsumer) processDelivery(ctx context.Context, delivery amqp.Delivery, handler MessageHandler) {
	retryCount := GetRetryCount(delivery)

	c.logger.Debug("Processing message",
		"messageId", delivery.MessageId,
		"correlationId", delivery.CorrelationId,
		"retryCount", retryCount,
	)

	err := handler(ctx, delivery)
	if err == nil {
		// Success - ACK
		if ackErr := delivery.Ack(false); ackErr != nil {
			c.logger.Error("Failed to ACK message", "error", ackErr, "messageId", delivery.MessageId)
		}
		return
	}

	// Handler returned error - determine retry or DLQ
	c.logger.Warn("Handler failed",
		"error", err,
		"messageId", delivery.MessageId,
		"retryCount", retryCount,
		"maxRetries", c.publisher.MaxRetries(),
	)

	maxRetries := c.publisher.MaxRetries()
	if retryCount < maxRetries {
		// Try to republish to retry queue first (before ACK)
		if pubErr := c.publishToRetryWithCount(ctx, retryCount, delivery); pubErr != nil {
			// Publish failed - NACK to requeue for redelivery
			c.logger.Error("Failed to publish to retry queue, requeueing",
				"error", pubErr,
				"retryIndex", retryCount,
				"messageId", delivery.MessageId,
			)
			if nackErr := delivery.Nack(false, true); nackErr != nil {
				c.logger.Error("Failed to NACK message for requeue", "error", nackErr)
			}
			return
		}

		// Publish succeeded - now ACK the original
		if ackErr := delivery.Ack(false); ackErr != nil {
			c.logger.Error("Failed to ACK after retry publish", "error", ackErr)
			// Message is in retry queue, duplicate may occur on redelivery
			return
		}

		c.logger.Info("Message queued for retry",
			"messageId", delivery.MessageId,
			"retryIndex", retryCount,
			"nextRetryCount", retryCount+1,
		)
	} else {
		// Max retries exhausted - NACK to DLQ
		c.logger.Error("Max retries exhausted, sending to DLQ",
			"messageId", delivery.MessageId,
			"retryCount", retryCount,
		)
		if nackErr := delivery.Nack(false, false); nackErr != nil {
			c.logger.Error("Failed to NACK message", "error", nackErr)
		}
	}
}

// publishToRetryWithCount publishes to retry queue with incremented retry count header
func (c *RabbitMQConsumer) publishToRetryWithCount(ctx context.Context, currentRetry int, delivery amqp.Delivery) error {
	// Create new headers with incremented retry count
	headers := make(amqp.Table)
	for k, v := range delivery.Headers {
		headers[k] = v
	}

	// Safe conversion: retry count is bounded by MaxRetries (typically < 10)
	nextRetry := currentRetry + 1
	headers[RetryCountHeader] = int32(nextRetry) //nolint:gosec // bounded by MaxRetries check

	// Use publisher's channel for retry publish
	return c.publisher.PublishToRetry(ctx, currentRetry, delivery.Body, delivery.CorrelationId, headers)
}

// GetRetryCount extracts the retry count from message headers
func GetRetryCount(delivery amqp.Delivery) int {
	if delivery.Headers == nil {
		return 0
	}

	val, ok := delivery.Headers[RetryCountHeader]
	if !ok {
		return 0
	}

	switch v := val.(type) {
	case int32:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}

// Connection returns the underlying AMQP connection for health checks
func (c *RabbitMQConsumer) Connection() *amqp.Connection {
	return c.conn
}

// Close stops consuming and closes the connection
func (c *RabbitMQConsumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true

	var errs []error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("channel: %v", err))
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("connection: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %w", ErrCloseFailed, errors.Join(errs...))
	}
	return nil
}
