# queue

RabbitMQ message publishing with retry support and dead letter queues.

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/queue"

publisher, err := queue.NewRabbitMQPublisher(cfg)
if err != nil {
    log.Fatal(err)
}
defer publisher.Close()

// Publish message
err = publisher.Publish(ctx, message)

// Retry failed message
err = publisher.PublishToRetry(ctx, retryIndex, body, correlationID)

// Send to dead letter queue
err = publisher.PublishToDLQ(ctx, body, correlationID)

// For health checks
conn := publisher.Connection()
```

## Features

- Automatic exchange and queue declaration
- Configurable retry delays with TTL-based routing
- Dead letter queue for permanent failures
- Connection accessor for health checks
- Thread-safe publishing

## Connection Ownership

The publisher owns and manages the RabbitMQ connection. `Connection()` returns
the underlying connection for read-only purposes (e.g., health checks).
**Do not call `conn.Close()` directly** - use `publisher.Close()` instead.

## Configuration

```go
cfg := config.RabbitMQConfig{
    Host:        "localhost",
    Port:        5672,
    User:        "guest",
    Password:    "guest",
    VHost:       "/",
    Exchange:    "messaging",
    Queue:       "contact_messages",
    RetryDelays: []time.Duration{5*time.Second, 30*time.Second, 2*time.Minute},
}
```
