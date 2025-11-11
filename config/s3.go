package config

// S3Config holds S3/MinIO storage configuration
type S3Config struct {
	Endpoint  string `validate:"required,url"`
	AccessKey string `validate:"required"`
	SecretKey string `validate:"required"`
	UseSSL    bool
}

// NewS3Config loads S3 configuration from environment variables
func NewS3Config() S3Config {
	return S3Config{
		Endpoint:  GetEnvRequired("S3_ENDPOINT"),
		AccessKey: GetEnvRequired("S3_ACCESS_KEY"),
		SecretKey: GetEnvRequired("S3_SECRET_KEY"),
		UseSSL:    GetEnvBool("S3_USE_SSL", false),
	}
}
