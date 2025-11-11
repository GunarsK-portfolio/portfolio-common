package config

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required,number,min=1,max=65535"`
	Password string // Optional, no validation
}

// LoadRedisConfig loads Redis configuration from environment variables
func LoadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     GetEnvRequired("REDIS_HOST"),
		Port:     GetEnvRequired("REDIS_PORT"),
		Password: GetEnv("REDIS_PASSWORD", ""),
	}
}
