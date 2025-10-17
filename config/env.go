package config

import (
	"log"
	"os"
	"strconv"
	"strings"
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
func GetEnvBool(key string, defaultValue bool) bool {
	val := GetEnv(key, "")
	if val == "" {
		return defaultValue
	}
	return strings.EqualFold(val, "true") || val == "1"
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
