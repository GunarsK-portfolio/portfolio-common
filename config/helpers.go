package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnv returns environment variable value or default if not set
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvRequired returns environment variable value or panics if not set
func GetEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

// GetEnvBool returns environment variable as boolean or default if not set
// (empty or whitespace-only counts as unset). Accepted values are
// case-insensitive true/false/1/0 with surrounding whitespace ignored; any
// other value panics at startup so a typo cannot silently flip a flag.
func GetEnvBool(key string, defaultValue bool) bool {
	val := strings.TrimSpace(GetEnv(key, ""))
	if val == "" {
		return defaultValue
	}
	b, ok := parseBool(val)
	if !ok {
		panic(fmt.Sprintf("Invalid boolean value for %s: %q (accepted: true, false, 1, 0)", key, val))
	}
	return b
}

// parseBool interprets the accepted boolean forms: case-insensitive
// true/false/1/0. ok is false for anything else. Deliberately narrower than
// strconv.ParseBool: the single-letter t/f forms are cryptic in env files
// and widen the accepted-typo surface.
func parseBool(val string) (value, ok bool) {
	switch strings.ToLower(val) {
	case "true", "1":
		return true, true
	case "false", "0":
		return false, true
	}
	return false, false
}

// GetEnvInt returns environment variable as integer or default if not set
func GetEnvInt(key string, defaultValue int) int {
	val := GetEnv(key, "")
	if val == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Warning: invalid integer value for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return intVal
}

// GetEnvInt64 returns environment variable as int64 or default if not set
func GetEnvInt64(key string, defaultValue int64) int64 {
	val := GetEnv(key, "")
	if val == "" {
		return defaultValue
	}
	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Printf("Warning: invalid int64 value for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return intVal
}

// GetEnvDuration returns environment variable as time.Duration or default if not set.
// Expected format: Go duration strings (e.g., "15m", "1h30m", "24h", "168h").
// See https://pkg.go.dev/time#ParseDuration for full format specification.
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	val := GetEnv(key, "")
	if val == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("Warning: invalid duration value for %s, using default %v", key, defaultValue)
		return defaultValue
	}
	return duration
}
