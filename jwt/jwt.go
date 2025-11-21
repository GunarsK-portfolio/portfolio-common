package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT token claims.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Service defines JWT token operations.
type Service interface {
	ValidateToken(tokenString string) (*Claims, error)
	GenerateAccessToken(userID int64, username string) (string, error)
	GenerateRefreshToken(userID int64, username string) (string, error)
	GetAccessExpiry() time.Duration
	GetRefreshExpiry() time.Duration
}

type service struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewService creates a new JWT Service instance.
// Returns error if secret is empty or less than 32 bytes for security.
func NewService(secret string, accessExpiry, refreshExpiry time.Duration) (Service, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 bytes")
	}
	return &service{
		secret:        secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}, nil
}

// NewValidatorOnly creates a JWT service for validation only (no token generation).
// Use this for services that only need to validate tokens, not generate them.
func NewValidatorOnly(secret string) (Service, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 bytes")
	}
	return &service{
		secret:        secret,
		accessExpiry:  0,
		refreshExpiry: 0,
	}, nil
}

func (s *service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *service) GenerateAccessToken(userID int64, username string) (string, error) {
	if s.accessExpiry == 0 {
		return "", errors.New("access token generation not configured")
	}
	return s.generateToken(userID, username, s.accessExpiry)
}

func (s *service) GenerateRefreshToken(userID int64, username string) (string, error) {
	if s.refreshExpiry == 0 {
		return "", errors.New("refresh token generation not configured")
	}
	return s.generateToken(userID, username, s.refreshExpiry)
}

func (s *service) generateToken(userID int64, username string, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *service) GetAccessExpiry() time.Duration {
	return s.accessExpiry
}

func (s *service) GetRefreshExpiry() time.Duration {
	return s.refreshExpiry
}

// GetTTL returns the remaining time-to-live for the token in seconds.
// Returns 0 if the token has expired or has no expiry.
func (c *Claims) GetTTL() int64 {
	if c.ExpiresAt == nil {
		return 0
	}
	ttl := time.Until(c.ExpiresAt.Time).Seconds()
	if ttl < 0 {
		return 0
	}
	return int64(ttl)
}
