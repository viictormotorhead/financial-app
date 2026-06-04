package investment

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	appinvestment "github.com/viictormotorhead/financial-app/internal/application/investment"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	apihandler "github.com/viictormotorhead/financial-app/internal/infra/api/handler"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/request"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
	apiresponse "github.com/viictormotorhead/financial-app/internal/infra/api/response"
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
)

func (h *InvestmentHandler) CreateValuation(c echo.Context) error {
	investmentID, err := parseInvestmentID(c)
	if err != nil {
		return err
	}

	var req request.CreateValuationRequest
	if err := c.Bind(&req); err != nil {
		return apiresponse.Error(http.StatusBadRequest, "request body must be valid JSON")
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

	return apiresponse.OK(c, http.StatusCreated, "valuation created successfully", map[string]response.ValuationResponse{
		"valuation": {
			ID:              output.HistoryID,
			InvestmentID:    output.InvestmentID,
			PreviousBalance: output.PreviousBalance,
			CurrentValue:    output.CurrentValue,
			Delta:           output.Delta,
			Balance:         output.Balance,
			Date:            output.Date,
		},
	})
}

func mapValuationError(err error) *echo.HTTPError {
	if httpErr := apihandler.HTTPErrorFromUseCase(err); httpErr != nil {
		return httpErr
	}
	switch {
	case errors.Is(err, appinvestment.ErrInvestmentNotFound):
		return apiresponse.Error(http.StatusNotFound, err.Error())
	case errors.Is(err, appinvestment.ErrValueUnchanged):
		return apiresponse.Error(http.StatusConflict, err.Error())
	default:
		return apiresponse.Error(http.StatusInternalServerError, "could not create valuation")
	}
}
