package config

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required,min=1,max=65535"`
	Password string // Optional, no validation
}

// NewRedisConfig loads Redis configuration from environment variables.
// It panics if required environment variables are missing or configuration is invalid.
func NewRedisConfig() RedisConfig {
	port, err := strconv.Atoi(GetEnvRequired("REDIS_PORT"))
	if err != nil {
		panic(fmt.Sprintf("Invalid REDIS_PORT: %v", err))
	}

	cfg := RedisConfig{
		Host:     GetEnvRequired("REDIS_HOST"),
		Port:     port,
		Password: GetEnv("REDIS_PASSWORD", ""),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid Redis configuration: %v", err))
	}

	return cfg
}
