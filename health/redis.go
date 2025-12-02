package health

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisChecker checks Redis connection status
type RedisChecker struct {
	client *redis.Client
}

// NewRedisChecker creates a new Redis health checker
func NewRedisChecker(client *redis.Client) Checker {
	return &RedisChecker{client: client}
}

// Name returns the name of this checker
func (c *RedisChecker) Name() string {
	return "redis"
}

// Check verifies the Redis connection by sending a PING command
func (c *RedisChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	if c.client == nil {
		return CheckResult{
			Status: StatusUnhealthy,
			Error:  "client is nil",
		}
	}

	if err := c.client.Ping(ctx).Err(); err != nil {
		return CheckResult{
			Status:  StatusUnhealthy,
			Latency: time.Since(start).String(),
			Error:   fmt.Sprintf("ping failed: %v", err),
		}
	}

	return CheckResult{
		Status:  StatusHealthy,
		Latency: time.Since(start).String(),
	}
}
