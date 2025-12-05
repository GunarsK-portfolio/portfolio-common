package config

import (
	"net/http"
	"strings"
)

// CookieConfig holds cookie configuration for authentication.
type CookieConfig struct {
	Domain      string // Cookie domain (e.g., ".gunarsk.com" for prod, "" for local)
	Secure      bool   // true for HTTPS only
	SameSite    http.SameSite
	Path        string // Path for access token cookie (typically "/")
	RefreshPath string // Path for refresh token cookie (must match refresh endpoint URL as seen by browser)
}

// NewCookieConfig loads cookie configuration from environment variables.
func NewCookieConfig() CookieConfig {
	sameSiteStr := strings.ToLower(GetEnv("COOKIE_SAMESITE", "Lax"))
	var sameSite http.SameSite
	switch sameSiteStr {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	case "lax", "":
		sameSite = http.SameSiteLaxMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	return CookieConfig{
		Domain:      GetEnv("COOKIE_DOMAIN", ""),
		Secure:      GetEnvBool("COOKIE_SECURE", false),
		SameSite:    sameSite,
		Path:        GetEnv("COOKIE_PATH", "/"),
		RefreshPath: GetEnv("COOKIE_REFRESH_PATH", "/"),
	}
}
