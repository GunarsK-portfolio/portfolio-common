# health

Dependency health checking with aggregated results.

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/health"

// Create aggregator with timeout
healthAgg := health.NewAggregator(3 * time.Second)

// Register checkers
healthAgg.Register(health.NewPostgresChecker(db))
healthAgg.Register(health.NewRabbitMQChecker(conn))
healthAgg.Register(health.NewRedisChecker(client))
healthAgg.Register(health.NewMinIOChecker(client, "bucket"))

// Use as Gin handler
router.GET("/health", healthAgg.Handler())
```

## Response Format

```json
{
  "status": "healthy",
  "checks": {
    "postgres": { "status": "healthy", "latency": "1.2ms" },
    "rabbitmq": { "status": "healthy", "latency": "0.3ms" }
  }
}
```

## HTTP Status Codes

- `200 OK` - All checks healthy
- `503 Service Unavailable` - Any check unhealthy or degraded

## Available Checkers

- `NewPostgresChecker(db *gorm.DB)` - PostgreSQL ping
- `NewRabbitMQChecker(conn *amqp.Connection)` - Connection status
- `NewRedisChecker(client *redis.Client)` - PING command
- `NewMinIOChecker(client *minio.Client, bucket string)` - Bucket check

## Status Types

- `healthy` - Check passed
- `degraded` - Partial failure (e.g., missing bucket)
- `unhealthy` - Check failed
