package config

import "testing"

func setRequiredRedisEnv(t *testing.T) {
	t.Helper()
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_PORT", "6379")
}

func TestNewRedisConfig_TLSDefaultsFalse(t *testing.T) {
	setRequiredRedisEnv(t)
	t.Setenv("REDIS_TLS", "")

	if cfg := NewRedisConfig(); cfg.TLS {
		t.Error("TLS = true, want false default")
	}
}

func TestNewRedisConfig_TLSEnabled(t *testing.T) {
	setRequiredRedisEnv(t)
	t.Setenv("REDIS_TLS", "true")

	if cfg := NewRedisConfig(); !cfg.TLS {
		t.Error("TLS = false, want true")
	}
}

func TestNewRedisConfig_InvalidTLSPanics(t *testing.T) {
	setRequiredRedisEnv(t)
	t.Setenv("REDIS_TLS", "yes")

	defer func() {
		if recover() == nil {
			t.Error("expected panic for invalid REDIS_TLS value")
		}
	}()
	NewRedisConfig()
}
