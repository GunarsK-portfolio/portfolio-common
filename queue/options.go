package queue

import "log/slog"

// PublisherOption configures a RabbitMQPublisher.
type PublisherOption func(*RabbitMQPublisher)

// WithPublisherLogger sets the logger used for reconnection events.
// Defaults to slog.Default().
func WithPublisherLogger(logger *slog.Logger) PublisherOption {
	return func(p *RabbitMQPublisher) {
		if logger != nil {
			p.logger = logger
		}
	}
}

// WithPublisherMetrics sets the metrics recorder for publish and reconnect
// events. Defaults to a no-op recorder.
func WithPublisherMetrics(m MetricsRecorder) PublisherOption {
	return func(p *RabbitMQPublisher) {
		if m != nil {
			p.metrics = m
		}
	}
}

// ConsumerOption configures a RabbitMQConsumer.
type ConsumerOption func(*RabbitMQConsumer)

// WithConsumerMetrics sets the metrics recorder for consume and reconnect
// events. Defaults to a no-op recorder.
func WithConsumerMetrics(m MetricsRecorder) ConsumerOption {
	return func(c *RabbitMQConsumer) {
		if m != nil {
			c.metrics = m
		}
	}
}
