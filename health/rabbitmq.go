package health

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQChecker checks RabbitMQ connection status
type RabbitMQChecker struct {
	conn *amqp.Connection
}

// NewRabbitMQChecker creates a new RabbitMQ health checker
func NewRabbitMQChecker(conn *amqp.Connection) Checker {
	return &RabbitMQChecker{conn: conn}
}

// Name returns the name of this checker
func (c *RabbitMQChecker) Name() string {
	return "rabbitmq"
}

// Check verifies the RabbitMQ connection is open
func (c *RabbitMQChecker) Check(_ context.Context) CheckResult {
	start := time.Now()

	if c.conn == nil {
		return CheckResult{
			Status:  StatusUnhealthy,
			Latency: time.Since(start).String(),
			Error:   "connection is nil",
		}
	}

	if c.conn.IsClosed() {
		return CheckResult{
			Status:  StatusUnhealthy,
			Latency: time.Since(start).String(),
			Error:   "connection is closed",
		}
	}

	return CheckResult{
		Status:  StatusHealthy,
		Latency: time.Since(start).String(),
	}
}
