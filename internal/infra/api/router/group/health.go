package group

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/viictormotorhead/financial-app/config"
)

const healthPath = "/health"

type HealthRoutes struct{}

func NewHealthRoutes(group *echo.Group) *HealthRoutes {
	base := config.Config().Server.BasePath
	healthRoutes := group.Group(base + healthPath)
	healthRoutes.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return &HealthRoutes{}
}
