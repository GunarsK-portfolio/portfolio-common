package config

import "time"

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret        string        `validate:"required,min=32"`
	AccessExpiry  time.Duration `validate:"gt=0"`
	RefreshExpiry time.Duration `validate:"gt=0"`
}

// LoadJWTConfig loads JWT configuration from environment variables
func LoadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:        GetEnvRequired("JWT_SECRET"),
		AccessExpiry:  GetEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		RefreshExpiry: GetEnvDuration("JWT_REFRESH_EXPIRY", 168*time.Hour),
	}
}
