package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "this-is-a-test-secret-with-32-bytes!" // 36 bytes

func TestNewService(t *testing.T) {
	tests := []struct {
		name          string
		secret        string
		accessExpiry  time.Duration
		refreshExpiry time.Duration
		wantErr       error
	}{
		{
			name:          "valid configuration",
			secret:        testSecret,
			accessExpiry:  15 * time.Minute,
			refreshExpiry: 168 * time.Hour,
			wantErr:       nil,
		},
		{
			name:          "secret too short",
			secret:        "short",
			accessExpiry:  15 * time.Minute,
			refreshExpiry: 168 * time.Hour,
			wantErr:       ErrSecretTooShort,
		},
		{
			name:          "empty secret",
			secret:        "",
			accessExpiry:  15 * time.Minute,
			refreshExpiry: 168 * time.Hour,
			wantErr:       ErrSecretTooShort,
		},
		{
			name:          "exactly 32 bytes",
			secret:        "12345678901234567890123456789012",
			accessExpiry:  15 * time.Minute,
			refreshExpiry: 168 * time.Hour,
			wantErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewService(tt.secret, tt.accessExpiry, tt.refreshExpiry)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if svc == nil {
					t.Error("expected service, got nil")
				}
			}
		})
	}
}

func TestNewValidatorOnly(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr error
	}{
		{
			name:    "valid secret",
			secret:  testSecret,
			wantErr: nil,
		},
		{
			name:    "secret too short",
			secret:  "short",
			wantErr: ErrSecretTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewValidatorOnly(tt.secret)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if svc == nil {
					t.Error("expected service, got nil")
				}
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	svc, err := NewService(testSecret, 15*time.Minute, 168*time.Hour)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	token, err := svc.GenerateAccessToken(123, "testuser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	// Validate the generated token
	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Errorf("failed to validate generated token: %v", err)
	}

	if claims.UserID != 123 {
		t.Errorf("expected UserID 123, got %d", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("expected Username testuser, got %s", claims.Username)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	svc, err := NewService(testSecret, 15*time.Minute, 168*time.Hour)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	token, err := svc.GenerateRefreshToken(456, "anotheruser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	// Validate the generated token
	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Errorf("failed to validate generated token: %v", err)
	}

	if claims.UserID != 456 {
		t.Errorf("expected UserID 456, got %d", claims.UserID)
	}
}

func TestValidatorOnlyCannotGenerateTokens(t *testing.T) {
	svc, err := NewValidatorOnly(testSecret)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	_, err = svc.GenerateAccessToken(123, "user")
	if err == nil {
		t.Error("expected error when generating access token with validator-only service")
	}

	_, err = svc.GenerateRefreshToken(123, "user")
	if err == nil {
		t.Error("expected error when generating refresh token with validator-only service")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	svc, err := NewService(testSecret, 15*time.Minute, 168*time.Hour)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	token, _ := svc.GenerateAccessToken(789, "validuser")
	claims, err := svc.ValidateToken(token)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if claims.UserID != 789 {
		t.Errorf("expected UserID 789, got %d", claims.UserID)
	}
	if claims.Username != "validuser" {
		t.Errorf("expected Username validuser, got %s", claims.Username)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	svc, err := NewService(testSecret, 1*time.Millisecond, 168*time.Hour)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	token, _ := svc.GenerateAccessToken(123, "user")
	time.Sleep(10 * time.Millisecond) // Wait for token to expire

	_, err = svc.ValidateToken(token)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	svc1, _ := NewService(testSecret, 15*time.Minute, 168*time.Hour)
	svc2, _ := NewService("different-secret-that-is-32-bytes!!", 15*time.Minute, 168*time.Hour)

	token, _ := svc1.GenerateAccessToken(123, "user")
	_, err := svc2.ValidateToken(token)

	if err == nil {
		t.Error("expected error for invalid signature")
	}
}

func TestValidateToken_MalformedToken(t *testing.T) {
	svc, _ := NewValidatorOnly(testSecret)

	tests := []struct {
		name  string
		token string
	}{
		{"empty string", ""},
		{"random string", "not-a-valid-token"},
		{"partial jwt", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
		{"invalid base64", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.###.xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ValidateToken(tt.token)
			if err == nil {
				t.Errorf("expected error for malformed token: %s", tt.token)
			}
		})
	}
}

func TestValidateToken_WrongSigningMethod(t *testing.T) {
	// Create a token with a different signing method (none)
	claims := &Claims{
		UserID:   123,
		Username: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	svc, _ := NewValidatorOnly(testSecret)
	_, err := svc.ValidateToken(tokenString)

	if err == nil {
		t.Error("expected error for wrong signing method")
	}
}

func TestGetTTL(t *testing.T) {
	tests := []struct {
		name     string
		claims   *Claims
		wantZero bool
	}{
		{
			name: "valid future expiry",
			claims: &Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				},
			},
			wantZero: false,
		},
		{
			name: "expired",
			claims: &Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				},
			},
			wantZero: true,
		},
		{
			name:     "nil expiry",
			claims:   &Claims{},
			wantZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttl := tt.claims.GetTTL()
			if tt.wantZero && ttl != 0 {
				t.Errorf("expected TTL 0, got %d", ttl)
			}
			if !tt.wantZero && ttl <= 0 {
				t.Errorf("expected positive TTL, got %d", ttl)
			}
		})
	}
}

func TestGetExpiry(t *testing.T) {
	accessExpiry := 15 * time.Minute
	refreshExpiry := 168 * time.Hour

	svc, err := NewService(testSecret, accessExpiry, refreshExpiry)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	if svc.GetAccessExpiry() != accessExpiry {
		t.Errorf("expected access expiry %v, got %v", accessExpiry, svc.GetAccessExpiry())
	}
	if svc.GetRefreshExpiry() != refreshExpiry {
		t.Errorf("expected refresh expiry %v, got %v", refreshExpiry, svc.GetRefreshExpiry())
	}

	// Validator-only should return 0
	validator, _ := NewValidatorOnly(testSecret)
	if validator.GetAccessExpiry() != 0 {
		t.Errorf("expected 0 access expiry for validator, got %v", validator.GetAccessExpiry())
	}
	if validator.GetRefreshExpiry() != 0 {
		t.Errorf("expected 0 refresh expiry for validator, got %v", validator.GetRefreshExpiry())
	}
}

func TestGenerateToken_InputValidation(t *testing.T) {
	svc, err := NewService(testSecret, 15*time.Minute, 168*time.Hour)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	tests := []struct {
		name     string
		userID   int64
		username string
		wantErr  error
	}{
		{
			name:     "valid inputs",
			userID:   123,
			username: "testuser",
			wantErr:  nil,
		},
		{
			name:     "zero user ID",
			userID:   0,
			username: "testuser",
			wantErr:  ErrInvalidUserID,
		},
		{
			name:     "negative user ID",
			userID:   -1,
			username: "testuser",
			wantErr:  ErrInvalidUserID,
		},
		{
			name:     "empty username",
			userID:   123,
			username: "",
			wantErr:  ErrEmptyUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" (access)", func(t *testing.T) {
			_, err := svc.GenerateAccessToken(tt.userID, tt.username)
			if err != tt.wantErr {
				t.Errorf("GenerateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		t.Run(tt.name+" (refresh)", func(t *testing.T) {
			_, err := svc.GenerateRefreshToken(tt.userID, tt.username)
			if err != tt.wantErr {
				t.Errorf("GenerateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
