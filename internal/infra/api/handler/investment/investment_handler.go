package investment

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/request"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
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
		return echo.NewHTTPError(http.StatusBadRequest, apivalidation.ErrorResponse{
			Message: "invalid request body",
			Errors:  []apivalidation.FieldError{{Field: "", Message: "request body must be valid JSON"}},
		})
	}

	if err := req.Validate(); err != nil {
		return apivalidation.HTTPError(err)
	}

	output, err := h.createInvestment.Create(c.Request().Context(), use_cases.CreateInvestmentCommand{
		Name:    req.Name,
		Balance: req.Balance,
		Tags:    req.Tags,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"message": "could not create investment",
		})
	}

	tags := output.Tags
	if tags == nil {
		tags = []string{}
	}

	return c.JSON(http.StatusCreated, response.InvestmentResponse{
		ID:        output.ID,
		Name:      output.Name,
		Balance:   output.Balance,
		Tags:      tags,
		CreatedAt: output.CreatedAt,
	})
}
