package config

// S3Config holds S3/MinIO storage configuration
// AccessKey and SecretKey are optional - if not provided, IAM role credentials will be used (AWS only)
type S3Config struct {
	Endpoint  string `validate:"required,url"`
	AccessKey string // Optional: required for MinIO, empty for AWS IAM role auth
	SecretKey string // Optional: required for MinIO, empty for AWS IAM role auth
	UseSSL    bool
}

// NewS3Config loads S3 configuration from environment variables
func NewS3Config() S3Config {
	cfg := S3Config{
		Endpoint:  GetEnvRequired("S3_ENDPOINT"),
		AccessKey: GetEnv("S3_ACCESS_KEY", ""), // Optional for IAM role authentication
		SecretKey: GetEnv("S3_SECRET_KEY", ""), // Optional for IAM role authentication
		UseSSL:    GetEnvBool("S3_USE_SSL", false),
	}

	// Only validate endpoint (access keys are optional for IAM role auth)
	if cfg.Endpoint == "" {
		panic("S3_ENDPOINT is required")
	}

	return cfg
}
