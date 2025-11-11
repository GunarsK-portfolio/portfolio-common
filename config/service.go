package config

// ServiceConfig holds service-level configuration (port, environment)
type ServiceConfig struct {
	Port        string `validate:"required,number,min=1,max=65535"`
	Environment string `validate:"oneof=development staging production"`
}

// LoadServiceConfig loads service configuration from environment variables
func LoadServiceConfig(defaultPort string) ServiceConfig {
	return ServiceConfig{
		Port:        GetEnv("PORT", defaultPort),
		Environment: GetEnv("ENVIRONMENT", "development"),
	}
}
