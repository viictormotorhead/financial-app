package investment

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	appinvestment "github.com/viictormotorhead/financial-app/internal/application/investment"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/request"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
)

func (h *InvestmentHandler) CreateValuation(c echo.Context) error {
	investmentID, err := parseInvestmentID(c)
	if err != nil {
		return err
	}

	var req request.CreateValuationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, apivalidation.ErrorResponse{
			Message: "invalid request body",
			Errors:  []apivalidation.FieldError{{Field: "", Message: "request body must be valid JSON"}},
		})
	}

	if err := req.Validate(); err != nil {
		return apivalidation.HTTPError(err)
	}

	output, err := h.createValuation.Create(c.Request().Context(), use_cases.CreateValuationCommand{
		InvestmentID: investmentID,
		CurrentValue: req.CurrentValue,
	})
	if err != nil {
		return mapValuationError(err)
	}

	return c.JSON(http.StatusCreated, response.ValuationResponse{
		ID:              output.HistoryID,
		InvestmentID:    output.InvestmentID,
		PreviousBalance: output.PreviousBalance,
		CurrentValue:    output.CurrentValue,
		Delta:           output.Delta,
		Balance:         output.Balance,
		Date:            output.Date,
	})
}

func mapValuationError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, appinvestment.ErrInvestmentNotFound):
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"message": err.Error()})
	case errors.Is(err, appinvestment.ErrValueUnchanged):
		return echo.NewHTTPError(http.StatusConflict, map[string]string{"message": err.Error()})
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"message": "could not create valuation",
		})
	}
}
