# queue

RabbitMQ message publishing and consuming with retry support and dead letter queues.

## Publisher Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/queue"

publisher, err := queue.NewRabbitMQPublisher(cfg)
if err != nil {
    log.Fatal(err)
}
defer publisher.Close()

// Publish message
err = publisher.Publish(ctx, message)

// Retry failed message (with headers for retry tracking)
err = publisher.PublishToRetry(ctx, retryIndex, body, correlationID, headers)

// Send to dead letter queue
err = publisher.PublishToDLQ(ctx, body, correlationID)

// For health checks
conn := publisher.Connection()
```

## Consumer Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/queue"

consumer, err := queue.NewRabbitMQConsumer(cfg, publisher, logger)
if err != nil {
    log.Fatal(err)
}
defer consumer.Close()

// Consume messages (blocks until context cancelled)
err = consumer.Consume(ctx, func(c context.Context, d amqp.Delivery) error {
    // Process message
    // Return nil to ACK, error to trigger retry logic
    return nil
})

// Get retry count from message headers
retryCount := queue.GetRetryCount(delivery)
```

## Features

- Automatic exchange and queue declaration
- Configurable retry delays with TTL-based routing
- Dead letter queue for permanent failures
- Connection accessor for health checks
- Thread-safe publishing and consuming
- Automatic retry count tracking via headers

## Retry Flow

1. Message fails processing → consumer calls `PublishToRetry()`
2. Message routes to retry queue with incremented `x-retry-count` header
3. After TTL expires, message returns to main queue
4. When `retryCount >= MaxRetries()`, message is NACKed → routes to DLQ

## Connection Ownership

Both publisher and consumer own their own RabbitMQ connections. `Connection()`
returns the underlying connection for read-only purposes (e.g., health checks).
**Do not call `conn.Close()` directly** - use `publisher.Close()` or
`consumer.Close()` instead.

## Configuration

```go
cfg := config.RabbitMQConfig{
    Host:        "localhost",
    Port:        5672,
    User:        "guest",
    Password:    "guest",
    TLS:         false,            // Set to true for amqps:// (Amazon MQ, production)
    Exchange:    "messaging",
    Queue:       "contact_messages",
    RetryDelays: []time.Duration{1*time.Minute, 5*time.Minute, 30*time.Minute},

    // Consumer-specific (optional)
    PrefetchCount: 1,              // QoS prefetch count
    ConsumerTag:   "my-consumer",  // Unique consumer identifier
}
```

## Environment Variables

| Variable | Required | Default | Description |
| -------- | -------- | ------- | ----------- |
| `RABBITMQ_HOST` | Yes | - | RabbitMQ hostname |
| `RABBITMQ_PORT` | Yes | - | RabbitMQ port |
| `RABBITMQ_USER` | Yes | - | Username |
| `RABBITMQ_PASSWORD` | Yes | - | Password |
| `RABBITMQ_TLS` | No | `false` | Use TLS (amqps://) connection |
| `RABBITMQ_EXCHANGE` | No | `contact_messages` | Exchange name |
| `RABBITMQ_QUEUE` | No | `contact_messages` | Queue name |
| `RABBITMQ_RETRY_DELAYS` | No | `1m,5m,30m,2h,12h` | Comma-separated durations |
| `RABBITMQ_PREFETCH_COUNT` | No | `1` | Consumer QoS prefetch |
| `RABBITMQ_CONSUMER_TAG` | No | `""` | Consumer identifier |

## Queue Infrastructure

Created automatically by `NewRabbitMQPublisher`:

- **Main queue** (`contact_messages`) - Primary message queue
- **Retry queues** (`contact_messages_retry_0`, `_1`, etc.) - TTL-based delays
- **DLQ** (`contact_messages_dlq`) - Permanent failures
- **Exchanges** - Direct exchanges for routing
