package investment

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	appinvestment "github.com/viictormotorhead/financial-app/internal/application/investment"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/request"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
	apiresponse "github.com/viictormotorhead/financial-app/internal/infra/api/response"
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
)

func (h *InvestmentHandler) CreateMovement(c echo.Context) error {
	investmentID, err := parseInvestmentID(c)
	if err != nil {
		return err
	}

	var req request.CreateMovementRequest
	if err := c.Bind(&req); err != nil {
		return apiresponse.Error(http.StatusBadRequest, "request body must be valid JSON")
	}

	if err := req.Validate(); err != nil {
		return apivalidation.HTTPError(err)
	}

	output, err := h.createMovement.Create(c.Request().Context(), use_cases.CreateMovementCommand{
		InvestmentID: investmentID,
		Type:         req.MovementType(),
		Amount:       req.Amount,
	})
	if err != nil {
		return mapMovementError(err)
	}

	return apiresponse.OK(c, http.StatusCreated, "movement created successfully", map[string]response.MovementResponse{
		"movement": {
			ID:           output.HistoryID,
			InvestmentID: output.InvestmentID,
			Type:         output.Type,
			Amount:       output.Amount,
			Balance:      output.Balance,
			Date:         output.Date,
		},
	})
}

func mapMovementError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, appinvestment.ErrInvestmentNotFound):
		return apiresponse.Error(http.StatusNotFound, err.Error())
	case errors.Is(err, appinvestment.ErrInsufficientBalance):
		return apiresponse.Error(http.StatusUnprocessableEntity, err.Error())
	default:
		return apiresponse.Error(http.StatusInternalServerError, "could not create movement")
	}
}

func parseInvestmentID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, apiresponse.Error(http.StatusBadRequest, "invalid investment id")
	}
	return uint(id), nil
}
