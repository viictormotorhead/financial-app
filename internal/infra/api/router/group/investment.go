package group

import (
	"github.com/labstack/echo/v4"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment"
)

const InvestmentPathV1 = "/v1/investments"

type InvestmentRoutes struct{}

func NewInvestmentRoutes(group *echo.Group, handler investment.InvestmentHandlerIF) *InvestmentRoutes {
	routes := group.Group(InvestmentPathV1)
	routes.POST("/", handler.Create)

	return &InvestmentRoutes{}
}
