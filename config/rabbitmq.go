package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// RabbitMQConfig holds RabbitMQ connection configuration
type RabbitMQConfig struct {
	Host        string          `validate:"required"`
	Port        string          `validate:"required,number,min=1,max=65535"`
	User        string          `validate:"required"`
	Password    string          `validate:"required"`
	Exchange    string          `validate:"required"`
	Queue       string          `validate:"required"`
	RetryDelays []time.Duration // Delays for retry queues (e.g., 5s, 30s, 5m, 30m, 2h)
}

// DefaultRetryDelays provides sensible defaults for retry delays
// Designed for email delivery: quick retry for transient issues, longer waits for outages
var DefaultRetryDelays = []time.Duration{
	1 * time.Minute,  // Transient network issues
	5 * time.Minute,  // Service temporarily unavailable
	30 * time.Minute, // Longer outage
	2 * time.Hour,    // Extended issue
	12 * time.Hour,   // Major outage, last retry before permanent failure
}

// NewRabbitMQConfig loads RabbitMQ configuration from environment variables
func NewRabbitMQConfig() RabbitMQConfig {
	cfg := RabbitMQConfig{
		Host:        GetEnvRequired("RABBITMQ_HOST"),
		Port:        GetEnvRequired("RABBITMQ_PORT"),
		User:        GetEnvRequired("RABBITMQ_USER"),
		Password:    GetEnvRequired("RABBITMQ_PASSWORD"),
		Exchange:    GetEnv("RABBITMQ_EXCHANGE", "contact_messages"),
		Queue:       GetEnv("RABBITMQ_QUEUE", "contact_messages"),
		RetryDelays: parseRetryDelays(GetEnv("RABBITMQ_RETRY_DELAYS", "")),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid RabbitMQ configuration: %v", err))
	}

	return cfg
}

// parseRetryDelays parses comma-separated duration strings (e.g., "5s,30s,5m,30m,2h")
func parseRetryDelays(s string) []time.Duration {
	if s == "" {
		return DefaultRetryDelays
	}

	parts := strings.Split(s, ",")
	delays := make([]time.Duration, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		d, err := time.ParseDuration(part)
		if err != nil {
			panic(fmt.Sprintf("Invalid retry delay %q: %v", part, err))
		}
		delays = append(delays, d)
	}

	if len(delays) == 0 {
		return DefaultRetryDelays
	}

	return delays
}

// URL returns the AMQP connection URL
func (c RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", c.User, c.Password, c.Host, c.Port)
}
