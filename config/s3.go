package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// S3Config holds S3/MinIO storage configuration
type S3Config struct {
	Endpoint  string `validate:"required,url"`
	AccessKey string `validate:"required"`
	SecretKey string `validate:"required"`
	UseSSL    bool
}

// NewS3Config loads S3 configuration from environment variables
func NewS3Config() S3Config {
	cfg := S3Config{
		Endpoint:  GetEnvRequired("S3_ENDPOINT"),
		AccessKey: GetEnvRequired("S3_ACCESS_KEY"),
		SecretKey: GetEnvRequired("S3_SECRET_KEY"),
		UseSSL:    GetEnvBool("S3_USE_SSL", false),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid S3 configuration: %v", err))
	}

	return cfg
}
