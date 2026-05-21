package investment

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/request"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
)

type InvestmentHandlerIF interface {
	Create(c echo.Context) error
}

type InvestmentHandler struct {
	createInvestment use_cases.CreateInvestmentUseCaseIF
}

func NewInvestmentHandler(createInvestment use_cases.CreateInvestmentUseCaseIF) InvestmentHandlerIF {
	return &InvestmentHandler{createInvestment: createInvestment}
}

func (h *InvestmentHandler) Create(c echo.Context) error {
	var req request.CreateInvestmentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	output, err := h.createInvestment.Create(c.Request().Context(), use_cases.CreateInvestmentCommand{
		Name:    req.Name,
		Balance: req.Balance,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, response.InvestmentResponse{
		ID:        output.ID,
		Name:      output.Name,
		Balance:   output.Balance,
		CreatedAt: output.CreatedAt,
	})
}
