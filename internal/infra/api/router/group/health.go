package group

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/viictormotorhead/financial-app/internal/infra/api/response"
	"github.com/viictormotorhead/financial-app/internal/infra/config"
)

const healthPath = "/health"

type HealthRoutes struct{}

func NewHealthRoutes(group *echo.Group) *HealthRoutes {
	base := config.Config().Server.BasePath
	healthRoutes := group.Group(base + healthPath)
	healthRoutes.GET("", func(c echo.Context) error {
		return response.OK(c, http.StatusOK, "service is healthy", nil)
	})

	return &HealthRoutes{}
}
