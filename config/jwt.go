package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret        string        `validate:"required,min=32"`
	AccessExpiry  time.Duration `validate:"gt=0"`
	RefreshExpiry time.Duration `validate:"gt=0"`
}

// NewJWTConfig loads JWT configuration from environment variables.
// Default values:
//   - JWT_ACCESS_EXPIRY: 15m (15 minutes)
//   - JWT_REFRESH_EXPIRY: 168h (7 days)
func NewJWTConfig() JWTConfig {
	cfg := JWTConfig{
		Secret:        GetEnvRequired("JWT_SECRET"),
		AccessExpiry:  GetEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		RefreshExpiry: GetEnvDuration("JWT_REFRESH_EXPIRY", 168*time.Hour),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid JWT configuration: %v", err))
	}

	return cfg
}
