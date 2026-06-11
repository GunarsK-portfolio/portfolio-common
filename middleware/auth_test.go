package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newClaimsTestContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	return c
}

func TestGetClaims_FullSet(t *testing.T) {
	c := newClaimsTestContext()
	c.Set(CtxKeyUserID, int64(42))
	c.Set(CtxKeyUsername, "kaladin")
	c.Set(CtxKeyDisplayName, "Kaladin Stormblessed")
	c.Set(CtxKeyScopes, map[string]string{"heroes": "edit"})

	claims, ok := GetClaims(c)
	if !ok {
		t.Fatal("GetClaims() ok = false, want true")
	}
	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Username != "kaladin" {
		t.Errorf("Username = %q, want %q", claims.Username, "kaladin")
	}
	if claims.DisplayName != "Kaladin Stormblessed" {
		t.Errorf("DisplayName = %q, want %q", claims.DisplayName, "Kaladin Stormblessed")
	}
	if claims.Scopes["heroes"] != "edit" {
		t.Errorf("Scopes = %v, want heroes:edit", claims.Scopes)
	}
}

func TestGetClaims_UserIDIsTheSentinel(t *testing.T) {
	c := newClaimsTestContext()
	c.Set(CtxKeyUserID, int64(1))

	claims, ok := GetClaims(c)
	if !ok {
		t.Fatal("GetClaims() ok = false, want true")
	}
	if claims.Username != "" || claims.DisplayName != "" || claims.Scopes != nil {
		t.Errorf("claims = %+v, want zero fields besides UserID", claims)
	}
}

func TestGetClaims_NotAuthenticated(t *testing.T) {
	tests := []struct {
		name string
		set  func(*gin.Context)
	}{
		{"nothing set", func(_ *gin.Context) {}},
		{"wrong user_id type", func(c *gin.Context) {
			c.Set(CtxKeyUserID, "1")
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClaimsTestContext()
			tt.set(c)
			if _, ok := GetClaims(c); ok {
				t.Error("GetClaims() ok = true, want false")
			}
		})
	}
}
