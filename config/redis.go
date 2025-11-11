package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required,number,min=1,max=65535"`
	Password string // Optional, no validation
}

// NewRedisConfig loads Redis configuration from environment variables
func NewRedisConfig() RedisConfig {
	cfg := RedisConfig{
		Host:     GetEnvRequired("REDIS_HOST"),
		Port:     GetEnvRequired("REDIS_PORT"),
		Password: GetEnv("REDIS_PASSWORD", ""),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid Redis configuration: %v", err))
	}

	return cfg
}
