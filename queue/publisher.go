package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/GunarsK-portfolio/portfolio-common/config"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Common errors returned by the queue package.
var (
	ErrConnectionFailed = errors.New("failed to connect to RabbitMQ")
	ErrChannelFailed    = errors.New("failed to open channel")
	ErrQueueSetupFailed = errors.New("failed to setup queue infrastructure")
	ErrMarshalFailed    = errors.New("failed to marshal message")
	ErrPublishFailed    = errors.New("failed to publish message")
	ErrRetryOutOfBounds = errors.New("retry index out of bounds")
	ErrCloseFailed      = errors.New("failed to close connection")
)

// Publisher defines the interface for message queue publishing with retry support
type Publisher interface {
	Publish(ctx context.Context, message interface{}) error
	PublishToRetry(ctx context.Context, retryIndex int, body []byte, correlationId string) error
	PublishToDLQ(ctx context.Context, body []byte, correlationId string) error
	MaxRetries() int
	Close() error
}

// RabbitMQPublisher implements Publisher for RabbitMQ.
// All publish methods are safe for concurrent use.
type RabbitMQPublisher struct {
	mu          sync.Mutex
	conn        *amqp.Connection
	channel     *amqp.Channel
	exchange    string
	queue       string
	retryQueues []string // Names of retry queues in order
}

// RetryQueues returns a copy of the retry queue names for use by consumers
func (p *RabbitMQPublisher) RetryQueues() []string {
	return append([]string(nil), p.retryQueues...)
}

// DLQName returns the dead letter queue name
func (p *RabbitMQPublisher) DLQName() string {
	return p.queue + "_dlq"
}

// declareExchangeAndQueue declares an exchange, queue, and binds them together
func declareExchangeAndQueue(ch *amqp.Channel, exchange, queue string, queueArgs amqp.Table) error {
	if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("%w: declare exchange %s: %v", ErrQueueSetupFailed, exchange, err)
	}

	if _, err := ch.QueueDeclare(queue, true, false, false, false, queueArgs); err != nil {
		return fmt.Errorf("%w: declare queue %s: %v", ErrQueueSetupFailed, queue, err)
	}

	if err := ch.QueueBind(queue, queue, exchange, false, nil); err != nil {
		return fmt.Errorf("%w: bind queue %s to exchange %s: %v", ErrQueueSetupFailed, queue, exchange, err)
	}

	return nil
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher with exchange, retry queues, and DLQ.
//
// The publisher is safe for concurrent use from multiple goroutines.
//
// Note: This publisher does not handle automatic reconnection. If the connection drops,
// callers should create a new publisher instance.
//
// Retry flow: Consumers must explicitly call PublishToRetry() to route failed messages through
// the retry chain. The main queue's dead-letter config routes directly to DLQ for unhandled
// failures (e.g., message rejected without calling PublishToRetry).
//
// If cfg.RetryDelays is empty, rejected messages route directly to the DLQ with no retry attempts.
func NewRabbitMQPublisher(cfg config.RabbitMQConfig) (*RabbitMQPublisher, error) {
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

	dlxName := cfg.Exchange + "_dlx"
	dlqName := cfg.Queue + "_dlq"

	// Declare dead letter exchange and queue (permanent failures)
	if err := declareExchangeAndQueue(ch, dlxName, dlqName, nil); err != nil {
		cleanup()
		return nil, err
	}

	// Declare retry queues with TTL (messages expire and route back to main queue)
	retryQueues := make([]string, len(cfg.RetryDelays))
	for i, delay := range cfg.RetryDelays {
		retryQueueName := fmt.Sprintf("%s_retry_%d", cfg.Queue, i)
		retryQueues[i] = retryQueueName

		retryArgs := amqp.Table{
			"x-message-ttl":             int64(delay.Milliseconds()),
			"x-dead-letter-exchange":    cfg.Exchange,
			"x-dead-letter-routing-key": cfg.Queue,
		}
		if err := declareExchangeAndQueue(ch, cfg.Exchange, retryQueueName, retryArgs); err != nil {
			cleanup()
			return nil, err
		}
	}

	// Declare main queue (failures route to first retry queue or DLQ if no retries left)
	mainQueueArgs := amqp.Table{
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": dlqName,
	}
	if err := declareExchangeAndQueue(ch, cfg.Exchange, cfg.Queue, mainQueueArgs); err != nil {
		cleanup()
		return nil, err
	}

	return &RabbitMQPublisher{
		conn:        conn,
		channel:     ch,
		exchange:    cfg.Exchange,
		queue:       cfg.Queue,
		retryQueues: retryQueues,
	}, nil
}

// publish is the internal helper for all publish operations.
// correlationId links related messages (e.g., original + retries).
func (p *RabbitMQPublisher) publish(ctx context.Context, exchange, routingKey string, body []byte, correlationId string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.channel.PublishWithContext(ctx, exchange, routingKey, false, false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "application/json",
			Body:          body,
			Timestamp:     time.Now(),
			MessageId:     uuid.NewString(),
			CorrelationId: correlationId,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPublishFailed, err)
	}
	return nil
}

// Publish sends a message to the main queue.
// Generates a new CorrelationId for tracking this message through retries.
func (p *RabbitMQPublisher) Publish(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMarshalFailed, err)
	}
	correlationId := uuid.NewString()
	return p.publish(ctx, p.exchange, p.queue, body, correlationId)
}

// PublishToRetry sends a message to a specific retry queue by index.
// Returns error if retryIndex is out of bounds (should send to DLQ instead).
// The correlationId should be preserved from the original message for tracing.
func (p *RabbitMQPublisher) PublishToRetry(ctx context.Context, retryIndex int, body []byte, correlationId string) error {
	maxRetries := p.MaxRetries()
	if maxRetries == 0 {
		return fmt.Errorf("%w: no retry queues configured", ErrRetryOutOfBounds)
	}
	if retryIndex < 0 || retryIndex >= maxRetries {
		return fmt.Errorf("%w: index %d, max %d", ErrRetryOutOfBounds, retryIndex, maxRetries-1)
	}
	return p.publish(ctx, p.exchange, p.retryQueues[retryIndex], body, correlationId)
}

// PublishToDLQ sends a message to the dead letter queue (permanent failure).
// The correlationId should be preserved from the original message for tracing.
func (p *RabbitMQPublisher) PublishToDLQ(ctx context.Context, body []byte, correlationId string) error {
	return p.publish(ctx, p.exchange+"_dlx", p.DLQName(), body, correlationId)
}

// MaxRetries returns the number of retry queues (attempts before DLQ)
func (p *RabbitMQPublisher) MaxRetries() int {
	return len(p.retryQueues)
}

// Close closes the channel and connection.
// Safe to call concurrently with Publish methods - will wait for in-flight publishes to complete.
func (p *RabbitMQPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var errs []error

	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("channel: %v", err))
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("connection: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %v", ErrCloseFailed, errs)
	}
	return nil
}
