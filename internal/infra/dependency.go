package infra

import (
	"github.com/labstack/echo/v4"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	"github.com/viictormotorhead/financial-app/internal/infra/api/middleware"
)

func NewEchoGroup(echoServer *echo.Echo, tokenParser appauth.TokenParser) *echo.Group {
	api := echoServer.Group("/api")
	api.Use(middleware.JWTAuth(tokenParser))
	return api
}
