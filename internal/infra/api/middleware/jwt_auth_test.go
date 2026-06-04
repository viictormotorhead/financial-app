package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
)

type stubTokenParser struct{}

func (stubTokenParser) Parse(string) (appauth.TokenClaims, error) {
	return appauth.TokenClaims{UserID: "user-1", Username: "victor"}, nil
}

func TestJWTAuth_PublicRoutes(t *testing.T) {
	e := echo.New()
	e.Use(JWTAuth(stubTokenParser{}))

	var reached bool
	e.POST(authPathLogin, func(c echo.Context) error {
		reached = true
		return c.NoContent(http.StatusOK)
	})
	e.POST(authPathRegister, func(c echo.Context) error {
		reached = true
		return c.NoContent(http.StatusOK)
	})
	e.GET("/api/v1/investments", func(c echo.Context) error {
		reached = true
		return c.NoContent(http.StatusOK)
	})

	for _, tc := range []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantReach  bool
	}{
		{name: "login without token", method: http.MethodPost, path: authPathLogin, wantStatus: http.StatusOK, wantReach: true},
		{name: "register without token", method: http.MethodPost, path: authPathRegister, wantStatus: http.StatusOK, wantReach: true},
		{name: "protected without token", method: http.MethodGet, path: "/api/v1/investments", wantStatus: http.StatusUnauthorized, wantReach: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reached = false
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status: got %d want %d", rec.Code, tc.wantStatus)
			}
			if reached != tc.wantReach {
				t.Fatalf("handler reached: got %v want %v", reached, tc.wantReach)
			}
		})
	}
}
