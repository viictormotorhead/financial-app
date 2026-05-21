package infra

import "github.com/labstack/echo/v4"

func NewEchoGroup(echoServer *echo.Echo) *echo.Group {
	return echoServer.Group("/api")
}
