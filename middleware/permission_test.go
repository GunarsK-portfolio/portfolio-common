package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		required string
		want     bool
	}{
		// Valid hierarchical comparisons
		{"delete >= delete", "delete", "delete", true},
		{"delete >= edit", "delete", "edit", true},
		{"delete >= read", "delete", "read", true},
		{"edit >= edit", "edit", "edit", true},
		{"edit >= read", "edit", "read", true},
		{"edit < delete", "edit", "delete", false},
		{"read >= read", "read", "read", true},
		{"read < edit", "read", "edit", false},
		{"read < delete", "read", "delete", false},
		{"none < read", "none", "read", false},
		{"none >= none", "none", "none", true},

		// Invalid/unknown levels - fail-safe behavior
		{"invalid user level", "invalid", "read", false},         // unknown user = 0 (none)
		{"invalid required level", "read", "invalid", false},     // unknown required = max (delete)
		{"empty string user level", "", "read", false},           // empty user = 0 (none)
		{"empty string required level", "read", "", false},       // empty required = max (delete)
		{"both unknown levels", "unknown1", "unknown2", false},   // 0 < max = denied
		{"delete with invalid required", "delete", "typo", true}, // 3 >= 3 (max)

		// Edge cases
		{"case sensitive - Read vs read", "Read", "read", false},
		{"whitespace - ' read ' vs read", " read ", "read", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPermission(tt.user, tt.required); got != tt.want {
				t.Errorf("HasPermission(%q, %q) = %v, want %v", tt.user, tt.required, got, tt.want)
			}
		})
	}
}

func TestRequirePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		scopes        interface{}
		scopesExists  bool
		resource      string
		requiredLevel string
		wantStatus    int
		wantError     string
		wantNext      bool
	}{
		{
			name:          "no scopes in context",
			scopesExists:  false,
			resource:      "projects",
			requiredLevel: "read",
			wantStatus:    http.StatusUnauthorized,
			wantError:     "unauthorized",
			wantNext:      false,
		},
		{
			name:          "invalid scopes format",
			scopes:        "not a map",
			scopesExists:  true,
			resource:      "projects",
			requiredLevel: "read",
			wantStatus:    http.StatusInternalServerError,
			wantError:     "invalid scopes format",
			wantNext:      false,
		},
		{
			name:          "sufficient permission - exact match",
			scopes:        map[string]string{"projects": "read"},
			scopesExists:  true,
			resource:      "projects",
			requiredLevel: "read",
			wantStatus:    http.StatusOK,
			wantNext:      true,
		},
		{
			name:          "sufficient permission - higher level",
			scopes:        map[string]string{"projects": "delete"},
			scopesExists:  true,
			resource:      "projects",
			requiredLevel: "read",
			wantStatus:    http.StatusOK,
			wantNext:      true,
		},
		{
			name:          "insufficient permission",
			scopes:        map[string]string{"projects": "read"},
			scopesExists:  true,
			resource:      "projects",
			requiredLevel: "edit",
			wantStatus:    http.StatusForbidden,
			wantError:     "insufficient permissions",
			wantNext:      false,
		},
		{
			name:          "resource not in scopes (empty string = none)",
			scopes:        map[string]string{"profile": "read"},
			scopesExists:  true,
			resource:      "projects",
			requiredLevel: "read",
			wantStatus:    http.StatusForbidden,
			wantError:     "insufficient permissions",
			wantNext:      false,
		},
		{
			name:          "none level cannot access read",
			scopes:        map[string]string{"users": "none"},
			scopesExists:  true,
			resource:      "users",
			requiredLevel: "read",
			wantStatus:    http.StatusForbidden,
			wantError:     "insufficient permissions",
			wantNext:      false,
		},
		{
			name:          "edit can access edit",
			scopes:        map[string]string{"skills": "edit"},
			scopesExists:  true,
			resource:      "skills",
			requiredLevel: "edit",
			wantStatus:    http.StatusOK,
			wantNext:      true,
		},
		{
			name:          "delete can access everything",
			scopes:        map[string]string{"files": "delete"},
			scopesExists:  true,
			resource:      "files",
			requiredLevel: "delete",
			wantStatus:    http.StatusOK,
			wantNext:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false

			router := gin.New()

			// Middleware to inject scopes into context
			router.Use(func(c *gin.Context) {
				if tt.scopesExists {
					c.Set("scopes", tt.scopes)
				}
				c.Next()
			})

			router.GET("/test",
				RequirePermission(tt.resource, tt.requiredLevel),
				func(c *gin.Context) {
					nextCalled = true
					c.Status(http.StatusOK)
				},
			)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			router.ServeHTTP(w, req)

			if nextCalled != tt.wantNext {
				t.Errorf("next handler called = %v, want %v", nextCalled, tt.wantNext)
			}

			if w.Code != tt.wantStatus {
				t.Errorf("status code = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantError != "" {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if response["error"] != tt.wantError {
					t.Errorf("error = %q, want %q", response["error"], tt.wantError)
				}
			}
		})
	}
}

func TestValidLevel(t *testing.T) {
	tests := []struct {
		level string
		want  bool
	}{
		{"none", true},
		{"read", true},
		{"edit", true},
		{"delete", true},
		{"invalid", false},
		{"", false},
		{"Read", false},  // case sensitive
		{" read", false}, // whitespace
		{"admin", false},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			if got := ValidLevel(tt.level); got != tt.want {
				t.Errorf("ValidLevel(%q) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestRequirePermission_PanicOnInvalidLevel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name  string
		level string
	}{
		{"typo reead", "reead"},
		{"empty string", ""},
		{"uppercase Read", "Read"},
		{"admin", "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("RequirePermission(resource, %q) did not panic", tt.level)
				}
			}()
			RequirePermission("projects", tt.level)
		})
	}
}

func TestLevelConstants(t *testing.T) {
	// Verify constants match expected values
	if LevelNone != "none" {
		t.Errorf("LevelNone = %q, want %q", LevelNone, "none")
	}
	if LevelRead != "read" {
		t.Errorf("LevelRead = %q, want %q", LevelRead, "read")
	}
	if LevelEdit != "edit" {
		t.Errorf("LevelEdit = %q, want %q", LevelEdit, "edit")
	}
	if LevelDelete != "delete" {
		t.Errorf("LevelDelete = %q, want %q", LevelDelete, "delete")
	}

	// Verify constants work with HasPermission
	if !HasPermission(LevelDelete, LevelRead) {
		t.Error("HasPermission(LevelDelete, LevelRead) should be true")
	}
	if HasPermission(LevelRead, LevelDelete) {
		t.Error("HasPermission(LevelRead, LevelDelete) should be false")
	}
}

func TestRequirePermission_ForbiddenResponseDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("scopes", map[string]string{"projects": "read"})
		c.Next()
	})
	router.GET("/test", RequirePermission("projects", "delete"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "insufficient permissions" {
		t.Errorf("error = %q, want %q", response["error"], "insufficient permissions")
	}
	if response["resource"] != "projects" {
		t.Errorf("resource = %q, want %q", response["resource"], "projects")
	}
	if response["required"] != "delete" {
		t.Errorf("required = %q, want %q", response["required"], "delete")
	}
	if response["have"] != "read" {
		t.Errorf("have = %q, want %q", response["have"], "read")
	}
}
