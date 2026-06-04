package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	"github.com/viictormotorhead/financial-app/internal/infra/api/response"
)

const (
	authPathRegister = "/api/v1/auth/register"
	authPathLogin    = "/api/v1/auth/login"
)

func JWTAuth(tokenParser appauth.TokenParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isPublicRoute(c) {
				return next(c)
			}

			token := bearerToken(c.Request().Header.Get("Authorization"))
			if token == "" {
				return response.Error(http.StatusUnauthorized, appauth.ErrUnauthorized.Error())
			}

			claims, err := tokenParser.Parse(token)
			if err != nil {
				return response.Error(http.StatusUnauthorized, appauth.ErrUnauthorized.Error())
			}

			ctx := appauth.WithUserID(c.Request().Context(), claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func isPublicRoute(c echo.Context) bool {
	if c.Request().Method == http.MethodOptions {
		return true
	}

	if c.Request().Method != http.MethodPost {
		return false
	}

	path := c.Request().URL.Path
	return path == authPathRegister || path == authPathLogin
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
