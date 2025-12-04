package config

import (
	"net/http"
)

// CookieConfig holds cookie configuration for authentication.
type CookieConfig struct {
	Domain   string // Cookie domain (e.g., ".gunarsk.com" for prod, "" for local)
	Secure   bool   // true for HTTPS only
	SameSite http.SameSite
	Path     string
}

// NewCookieConfig loads cookie configuration from environment variables.
func NewCookieConfig() CookieConfig {
	sameSiteStr := GetEnv("COOKIE_SAMESITE", "Lax")
	var sameSite http.SameSite
	switch sameSiteStr {
	case "Strict":
		sameSite = http.SameSiteStrictMode
	case "None":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	return CookieConfig{
		Domain:   GetEnv("COOKIE_DOMAIN", ""),
		Secure:   GetEnvBool("COOKIE_SECURE", false),
		SameSite: sameSite,
		Path:     GetEnv("COOKIE_PATH", "/"),
	}
}
