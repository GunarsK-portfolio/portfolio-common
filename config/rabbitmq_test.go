package config

import (
	"testing"
	"time"
)

// =============================================================================
// URL Tests
// =============================================================================

func TestRabbitMQConfig_URL(t *testing.T) {
	tests := []struct {
		name string
		cfg  RabbitMQConfig
		want string
	}{
		{
			name: "standard config",
			cfg: RabbitMQConfig{
				Host:     "localhost",
				Port:     5672,
				User:     "guest",
				Password: "guest",
			},
			want: "amqp://guest:guest@localhost:5672/",
		},
		{
			name: "special characters in password",
			cfg: RabbitMQConfig{
				Host:     "rabbitmq",
				Port:     5672,
				User:     "admin",
				Password: "p@ss:word/123",
			},
			want: "amqp://admin:p%40ss%3Aword%2F123@rabbitmq:5672/",
		},
		{
			name: "custom port",
			cfg: RabbitMQConfig{
				Host:     "192.168.1.100",
				Port:     15672,
				User:     "user",
				Password: "pass",
			},
			want: "amqp://user:pass@192.168.1.100:15672/",
		},
		{
			name: "TLS enabled uses amqps scheme",
			cfg: RabbitMQConfig{
				Host:     "b-xxx.mq.eu-west-1.on.aws",
				Port:     5671,
				User:     "rabbitmq",
				Password: "secret123",
				TLS:      true,
			},
			want: "amqps://rabbitmq:secret123@b-xxx.mq.eu-west-1.on.aws:5671/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.URL()
			if got != tt.want {
				t.Errorf("URL() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// DefaultRetryDelays Tests
// =============================================================================

func TestDefaultRetryDelays(t *testing.T) {
	delays := DefaultRetryDelays()

	// Should return 5 default delays
	if len(delays) != 5 {
		t.Errorf("DefaultRetryDelays() returned %d delays, want 5", len(delays))
	}

	// Verify expected values
	expected := []time.Duration{
		1 * time.Minute,
		5 * time.Minute,
		30 * time.Minute,
		2 * time.Hour,
		12 * time.Hour,
	}

	for i, want := range expected {
		if delays[i] != want {
			t.Errorf("DefaultRetryDelays()[%d] = %v, want %v", i, delays[i], want)
		}
	}

	// Verify it returns a copy (modifying returned slice shouldn't affect future calls)
	delays[0] = 999 * time.Hour
	newDelays := DefaultRetryDelays()
	if newDelays[0] == 999*time.Hour {
		t.Error("DefaultRetryDelays() should return a copy, not the original slice")
	}
}

// =============================================================================
// parseRetryDelays Tests
// =============================================================================

func TestParseRetryDelays(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantVals []time.Duration
	}{
		{
			name:    "empty string returns defaults",
			input:   "",
			wantLen: 5,
		},
		{
			name:     "single duration",
			input:    "30s",
			wantLen:  1,
			wantVals: []time.Duration{30 * time.Second},
		},
		{
			name:     "multiple durations",
			input:    "1m,5m,30m",
			wantLen:  3,
			wantVals: []time.Duration{1 * time.Minute, 5 * time.Minute, 30 * time.Minute},
		},
		{
			name:     "with spaces",
			input:    "1m, 5m, 30m",
			wantLen:  3,
			wantVals: []time.Duration{1 * time.Minute, 5 * time.Minute, 30 * time.Minute},
		},
		{
			name:     "hours",
			input:    "1h,2h,12h",
			wantLen:  3,
			wantVals: []time.Duration{1 * time.Hour, 2 * time.Hour, 12 * time.Hour},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delays := parseRetryDelays(tt.input)

			if len(delays) != tt.wantLen {
				t.Errorf("parseRetryDelays(%q) returned %d delays, want %d", tt.input, len(delays), tt.wantLen)
				return
			}

			if tt.wantVals != nil {
				for i, want := range tt.wantVals {
					if delays[i] != want {
						t.Errorf("parseRetryDelays(%q)[%d] = %v, want %v", tt.input, i, delays[i], want)
					}
				}
			}
		})
	}
}

func TestParseRetryDelays_Panics(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid duration format",
			input: "invalid",
		},
		{
			name:  "negative duration",
			input: "-5m",
		},
		{
			name:  "zero duration",
			input: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("parseRetryDelays(%q) should panic", tt.input)
				}
			}()

			parseRetryDelays(tt.input)
		})
	}
}

// =============================================================================
// Consumer Settings Tests
// =============================================================================

func TestRabbitMQConfig_ConsumerFields(t *testing.T) {
	cfg := RabbitMQConfig{
		Host:          "localhost",
		Port:          5672,
		User:          "guest",
		Password:      "guest",
		Exchange:      "test_exchange",
		Queue:         "test_queue",
		PrefetchCount: 10,
		ConsumerTag:   "my-consumer",
	}

	if cfg.PrefetchCount != 10 {
		t.Errorf("PrefetchCount = %d, want 10", cfg.PrefetchCount)
	}
	if cfg.ConsumerTag != "my-consumer" {
		t.Errorf("ConsumerTag = %q, want %q", cfg.ConsumerTag, "my-consumer")
	}
}

func TestRabbitMQConfig_ConsumerFieldsDefault(t *testing.T) {
	// When not set, consumer fields should be zero values
	cfg := RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		User:     "guest",
		Password: "guest",
		Exchange: "test_exchange",
		Queue:    "test_queue",
	}

	if cfg.PrefetchCount != 0 {
		t.Errorf("PrefetchCount default = %d, want 0", cfg.PrefetchCount)
	}
	if cfg.ConsumerTag != "" {
		t.Errorf("ConsumerTag default = %q, want empty string", cfg.ConsumerTag)
	}
}
