package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ServiceConfig holds service-level configuration (port, environment, CORS).
// Valid environment values: "development", "staging", "production"
type ServiceConfig struct {
	Port           string   `validate:"required,number,min=1,max=65535"`
	Environment    string   `validate:"oneof=development staging production"`
	AllowedOrigins []string `validate:"required,min=1"`
}

// NewServiceConfig loads service configuration from environment variables
func NewServiceConfig(defaultPort string) ServiceConfig {
	// Parse allowed origins from comma-separated string
	// NO DEFAULT - CORS must be explicitly configured
	allowedOriginsStr := GetEnvRequired("ALLOWED_ORIGINS")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	cfg := ServiceConfig{
		Port:           GetEnv("PORT", defaultPort),
		Environment:    GetEnv("ENVIRONMENT", "development"),
		AllowedOrigins: allowedOrigins,
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid service configuration: %v", err))
	}

	return cfg
}
