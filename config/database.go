package config

// DatabaseConfig holds PostgreSQL database configuration
type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required,number,min=1,max=65535"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
}

// NewDatabaseConfig loads database configuration from environment variables
func NewDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     GetEnvRequired("DB_HOST"),
		Port:     GetEnvRequired("DB_PORT"),
		User:     GetEnvRequired("DB_USER"),
		Password: GetEnvRequired("DB_PASSWORD"),
		Name:     GetEnvRequired("DB_NAME"),
	}
}
