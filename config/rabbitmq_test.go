package config

import (
	"strings"
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
// WithDefaults Tests
// =============================================================================

func TestWithDefaults_FillsZeroValues(t *testing.T) {
	cfg := RabbitMQConfig{}.WithDefaults()

	if cfg.Heartbeat != DefaultHeartbeat {
		t.Errorf("Heartbeat = %v, want %v", cfg.Heartbeat, DefaultHeartbeat)
	}
	if cfg.ReconnectInitialDelay != DefaultReconnectInitialDelay {
		t.Errorf("ReconnectInitialDelay = %v, want %v", cfg.ReconnectInitialDelay, DefaultReconnectInitialDelay)
	}
	if cfg.ReconnectMaxDelay != DefaultReconnectMaxDelay {
		t.Errorf("ReconnectMaxDelay = %v, want %v", cfg.ReconnectMaxDelay, DefaultReconnectMaxDelay)
	}
	if cfg.PrefetchCount != 1 {
		t.Errorf("PrefetchCount = %d, want 1", cfg.PrefetchCount)
	}
	if cfg.ConsumerConcurrency != 1 {
		t.Errorf("ConsumerConcurrency = %d, want 1", cfg.ConsumerConcurrency)
	}
}

func TestWithDefaults_KeepsExplicitValues(t *testing.T) {
	cfg := RabbitMQConfig{
		Heartbeat:             20 * time.Second,
		ReconnectInitialDelay: 2 * time.Second,
		ReconnectMaxDelay:     time.Minute,
		PrefetchCount:         8,
		ConsumerConcurrency:   4,
	}.WithDefaults()

	if cfg.Heartbeat != 20*time.Second {
		t.Errorf("Heartbeat = %v, want 20s", cfg.Heartbeat)
	}
	if cfg.ReconnectInitialDelay != 2*time.Second {
		t.Errorf("ReconnectInitialDelay = %v, want 2s", cfg.ReconnectInitialDelay)
	}
	if cfg.ReconnectMaxDelay != time.Minute {
		t.Errorf("ReconnectMaxDelay = %v, want 1m", cfg.ReconnectMaxDelay)
	}
	if cfg.PrefetchCount != 8 {
		t.Errorf("PrefetchCount = %d, want 8", cfg.PrefetchCount)
	}
	if cfg.ConsumerConcurrency != 4 {
		t.Errorf("ConsumerConcurrency = %d, want 4", cfg.ConsumerConcurrency)
	}
}

// =============================================================================
// Prefixed Environment Tests
// =============================================================================

func setBaseRabbitMQEnv(t *testing.T) {
	t.Helper()
	t.Setenv("RABBITMQ_HOST", "base-host")
	t.Setenv("RABBITMQ_PORT", "5672")
	t.Setenv("RABBITMQ_USER", "base-user")
	t.Setenv("RABBITMQ_PASSWORD", "base-pass")
	t.Setenv("RABBITMQ_QUEUE", "base_queue")
	t.Setenv("RABBITMQ_EXCHANGE", "base_exchange")
}

func TestNewRabbitMQConfigWithPrefix_EmptyPrefixUsesBaseVars(t *testing.T) {
	setBaseRabbitMQEnv(t)

	cfg := NewRabbitMQConfigWithPrefix("")

	if cfg.Host != "base-host" {
		t.Errorf("Host = %q, want %q", cfg.Host, "base-host")
	}
	if cfg.Queue != "base_queue" {
		t.Errorf("Queue = %q, want %q", cfg.Queue, "base_queue")
	}
}

func TestNewRabbitMQConfigWithPrefix_PrefixedOverridesBase(t *testing.T) {
	setBaseRabbitMQEnv(t)
	t.Setenv("AI_RABBITMQ_QUEUE", "ai_jobs")
	t.Setenv("AI_RABBITMQ_EXCHANGE", "ai_exchange")
	t.Setenv("AI_RABBITMQ_RETRY_DELAYS", "5s,15s")
	t.Setenv("AI_RABBITMQ_PREFETCH_COUNT", "2")

	cfg := NewRabbitMQConfigWithPrefix("AI_")

	// Prefixed values win.
	if cfg.Queue != "ai_jobs" {
		t.Errorf("Queue = %q, want %q", cfg.Queue, "ai_jobs")
	}
	if cfg.Exchange != "ai_exchange" {
		t.Errorf("Exchange = %q, want %q", cfg.Exchange, "ai_exchange")
	}
	if len(cfg.RetryDelays) != 2 || cfg.RetryDelays[0] != 5*time.Second {
		t.Errorf("RetryDelays = %v, want [5s 15s]", cfg.RetryDelays)
	}
	if cfg.PrefetchCount != 2 {
		t.Errorf("PrefetchCount = %d, want 2", cfg.PrefetchCount)
	}

	// Un-prefixed values are the fallback for shared settings.
	if cfg.Host != "base-host" {
		t.Errorf("Host = %q, want fallback %q", cfg.Host, "base-host")
	}
	if cfg.User != "base-user" {
		t.Errorf("User = %q, want fallback %q", cfg.User, "base-user")
	}
}

func TestNewRabbitMQConfig_Defaults(t *testing.T) {
	setBaseRabbitMQEnv(t)

	cfg := NewRabbitMQConfig()

	if cfg.Heartbeat != DefaultHeartbeat {
		t.Errorf("Heartbeat = %v, want %v", cfg.Heartbeat, DefaultHeartbeat)
	}
	if cfg.PublisherConfirms {
		t.Error("PublisherConfirms should default to false")
	}
	if cfg.DisableReconnect {
		t.Error("DisableReconnect should default to false (reconnect enabled)")
	}
	if cfg.ReconnectMaxAttempts != 0 {
		t.Errorf("ReconnectMaxAttempts = %d, want 0 (unlimited)", cfg.ReconnectMaxAttempts)
	}
	if cfg.ReconnectInitialDelay != DefaultReconnectInitialDelay {
		t.Errorf("ReconnectInitialDelay = %v, want %v", cfg.ReconnectInitialDelay, DefaultReconnectInitialDelay)
	}
	if cfg.ReconnectMaxDelay != DefaultReconnectMaxDelay {
		t.Errorf("ReconnectMaxDelay = %v, want %v", cfg.ReconnectMaxDelay, DefaultReconnectMaxDelay)
	}
	if cfg.RetryJitter != 0 {
		t.Errorf("RetryJitter = %v, want 0", cfg.RetryJitter)
	}
	if cfg.ConsumerConcurrency != 1 {
		t.Errorf("ConsumerConcurrency = %d, want 1", cfg.ConsumerConcurrency)
	}
}

func TestNewRabbitMQConfig_NewFieldsFromEnv(t *testing.T) {
	setBaseRabbitMQEnv(t)
	t.Setenv("RABBITMQ_HEARTBEAT", "20s")
	t.Setenv("RABBITMQ_PUBLISHER_CONFIRMS", "true")
	t.Setenv("RABBITMQ_RETRY_JITTER", "0.25")
	t.Setenv("RABBITMQ_RECONNECT", "false")
	t.Setenv("RABBITMQ_RECONNECT_MAX_ATTEMPTS", "10")
	t.Setenv("RABBITMQ_RECONNECT_INITIAL_DELAY", "2s")
	t.Setenv("RABBITMQ_RECONNECT_MAX_DELAY", "1m")
	t.Setenv("RABBITMQ_CONSUMER_CONCURRENCY", "4")

	cfg := NewRabbitMQConfig()

	if cfg.Heartbeat != 20*time.Second {
		t.Errorf("Heartbeat = %v, want 20s", cfg.Heartbeat)
	}
	if !cfg.PublisherConfirms {
		t.Error("PublisherConfirms should be true")
	}
	if cfg.RetryJitter != 0.25 {
		t.Errorf("RetryJitter = %v, want 0.25", cfg.RetryJitter)
	}
	if !cfg.DisableReconnect {
		t.Error("DisableReconnect should be true when RABBITMQ_RECONNECT=false")
	}
	if cfg.ReconnectMaxAttempts != 10 {
		t.Errorf("ReconnectMaxAttempts = %d, want 10", cfg.ReconnectMaxAttempts)
	}
	if cfg.ReconnectInitialDelay != 2*time.Second {
		t.Errorf("ReconnectInitialDelay = %v, want 2s", cfg.ReconnectInitialDelay)
	}
	if cfg.ReconnectMaxDelay != time.Minute {
		t.Errorf("ReconnectMaxDelay = %v, want 1m", cfg.ReconnectMaxDelay)
	}
	if cfg.ConsumerConcurrency != 4 {
		t.Errorf("ConsumerConcurrency = %d, want 4", cfg.ConsumerConcurrency)
	}
}

func TestNewRabbitMQConfig_BoolParsing(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		wantConfirms bool
	}{
		{"lowercase true", "true", true},
		{"uppercase with trailing space", "TRUE ", true},
		{"one", "1", true},
		{"mixed case false", "False", false},
		{"zero", "0", false},
		{"whitespace only keeps default", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setBaseRabbitMQEnv(t)
			t.Setenv("RABBITMQ_PUBLISHER_CONFIRMS", tt.value)

			cfg := NewRabbitMQConfig()

			if cfg.PublisherConfirms != tt.wantConfirms {
				t.Errorf("PublisherConfirms = %v for %q, want %v", cfg.PublisherConfirms, tt.value, tt.wantConfirms)
			}
		})
	}
}

func TestNewRabbitMQConfig_InvalidBoolPanics(t *testing.T) {
	tests := []struct {
		name        string
		envVar      string
		prefix      string
		wantInPanic string
	}{
		{
			name:        "yes is rejected naming the un-prefixed variable",
			envVar:      "RABBITMQ_RECONNECT",
			prefix:      "",
			wantInPanic: "RABBITMQ_RECONNECT",
		},
		{
			name:        "prefixed variable named when it resolved",
			envVar:      "AI_RABBITMQ_RECONNECT",
			prefix:      "AI_",
			wantInPanic: "AI_RABBITMQ_RECONNECT",
		},
		{
			name:        "un-prefixed fallback named when it resolved",
			envVar:      "RABBITMQ_RECONNECT",
			prefix:      "AI_",
			wantInPanic: "RABBITMQ_RECONNECT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setBaseRabbitMQEnv(t)
			t.Setenv(tt.envVar, "yes")

			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("NewRabbitMQConfigWithPrefix should panic for malformed bool")
				}
				msg, ok := r.(string)
				if !ok || !strings.Contains(msg, tt.wantInPanic) {
					t.Errorf("panic = %v, want it to name %s", r, tt.wantInPanic)
				}
			}()

			NewRabbitMQConfigWithPrefix(tt.prefix)
		})
	}
}

func TestNewRabbitMQConfig_InvalidJitterPanics(t *testing.T) {
	setBaseRabbitMQEnv(t)
	t.Setenv("RABBITMQ_RETRY_JITTER", "1.5")

	defer func() {
		if r := recover(); r == nil {
			t.Error("NewRabbitMQConfig should panic for jitter > 1")
		}
	}()

	NewRabbitMQConfig()
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
