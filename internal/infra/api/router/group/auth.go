package group

import (
	"github.com/labstack/echo/v4"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/auth"
)

const AuthPathV1 = "/v1/auth"

type AuthRoutes struct{}

func NewAuthRoutes(group *echo.Group, handler auth.AuthHandlerIF) *AuthRoutes {
	routes := group.Group(AuthPathV1)
	routes.POST("/register", handler.Register)
	routes.POST("/login", handler.Login)

	return &AuthRoutes{}
}
