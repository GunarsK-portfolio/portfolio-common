package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// RabbitMQConfig holds RabbitMQ connection configuration
type RabbitMQConfig struct {
	Host        string `validate:"required"`
	Port        int    `validate:"required,min=1,max=65535"`
	User        string `validate:"required"`
	Password    string `validate:"required"`
	Exchange    string `validate:"required"`
	Queue       string `validate:"required"`
	TLS         bool
	RetryDelays []time.Duration // Delays for retry queues (e.g., 5s, 30s, 5m, 30m, 2h)

	// Consumer-specific settings (optional, only used by consumers)
	PrefetchCount int    `validate:"omitempty,min=1"` // Number of messages to prefetch (QoS), defaults to 1
	ConsumerTag   string // Unique identifier for this consumer
}

// defaultRetryDelays provides sensible defaults for retry delays
// Designed for email delivery: quick retry for transient issues, longer waits for outages
var defaultRetryDelays = []time.Duration{
	1 * time.Minute,  // Transient network issues
	5 * time.Minute,  // Service temporarily unavailable
	30 * time.Minute, // Longer outage
	2 * time.Hour,    // Extended issue
	12 * time.Hour,   // Major outage, last retry before permanent failure
}

// DefaultRetryDelays returns a copy of the default retry delays
func DefaultRetryDelays() []time.Duration {
	return append([]time.Duration(nil), defaultRetryDelays...)
}

// NewRabbitMQConfig loads RabbitMQ configuration from environment variables.
// It panics if required environment variables are missing or configuration is invalid.
func NewRabbitMQConfig() RabbitMQConfig {
	port, err := strconv.Atoi(GetEnvRequired("RABBITMQ_PORT"))
	if err != nil {
		panic(fmt.Sprintf("Invalid RABBITMQ_PORT: %v", err))
	}

	cfg := RabbitMQConfig{
		Host:          GetEnvRequired("RABBITMQ_HOST"),
		Port:          port,
		User:          GetEnvRequired("RABBITMQ_USER"),
		Password:      GetEnvRequired("RABBITMQ_PASSWORD"),
		Exchange:      GetEnv("RABBITMQ_EXCHANGE", "contact_messages"),
		Queue:         GetEnv("RABBITMQ_QUEUE", "contact_messages"),
		TLS:           GetEnvBool("RABBITMQ_TLS", false),
		RetryDelays:   parseRetryDelays(GetEnv("RABBITMQ_RETRY_DELAYS", "")),
		PrefetchCount: GetEnvInt("RABBITMQ_PREFETCH_COUNT", 1),
		ConsumerTag:   GetEnv("RABBITMQ_CONSUMER_TAG", ""),
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
		return DefaultRetryDelays()
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
		if d <= 0 {
			panic(fmt.Sprintf("Retry delay must be positive, got %q", part))
		}
		delays = append(delays, d)
	}

	if len(delays) == 0 {
		return DefaultRetryDelays()
	}

	return delays
}

// URL returns the AMQP connection URL with properly encoded credentials.
// Uses amqps:// scheme when TLS is enabled, amqp:// otherwise.
func (c RabbitMQConfig) URL() string {
	scheme := "amqp"
	if c.TLS {
		scheme = "amqps"
	}
	u := &url.URL{
		Scheme: scheme,
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   "/",
	}
	return u.String()
}
