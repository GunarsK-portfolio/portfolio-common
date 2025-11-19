package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// DatabaseConfig holds PostgreSQL database configuration
type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required,number,min=1,max=65535"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
	SSLMode  string // optional, defaults to "disable" for local, set to "require" for AWS RDS
}

// NewDatabaseConfig loads database configuration from environment variables
func NewDatabaseConfig() DatabaseConfig {
	cfg := DatabaseConfig{
		Host:     GetEnvRequired("DB_HOST"),
		Port:     GetEnvRequired("DB_PORT"),
		User:     GetEnvRequired("DB_USER"),
		Password: GetEnvRequired("DB_PASSWORD"),
		Name:     GetEnvRequired("DB_NAME"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid database configuration: %v", err))
	}

	return cfg
}
