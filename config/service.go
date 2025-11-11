package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ServiceConfig holds service-level configuration (port, environment).
// Valid environment values: "development", "staging", "production"
type ServiceConfig struct {
	Port        string `validate:"required,number,min=1,max=65535"`
	Environment string `validate:"oneof=development staging production"`
}

// NewServiceConfig loads service configuration from environment variables
func NewServiceConfig(defaultPort string) ServiceConfig {
	cfg := ServiceConfig{
		Port:        GetEnv("PORT", defaultPort),
		Environment: GetEnv("ENVIRONMENT", "development"),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid service configuration: %v", err))
	}

	return cfg
}
