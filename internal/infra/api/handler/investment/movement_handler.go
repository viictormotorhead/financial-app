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
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
)

func (h *InvestmentHandler) CreateMovement(c echo.Context) error {
	investmentID, err := parseInvestmentID(c)
	if err != nil {
		return err
	}

	var req request.CreateMovementRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, apivalidation.ErrorResponse{
			Message: "invalid request body",
			Errors:  []apivalidation.FieldError{{Field: "", Message: "request body must be valid JSON"}},
		})
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

	return c.JSON(http.StatusCreated, response.MovementResponse{
		ID:           output.HistoryID,
		InvestmentID: output.InvestmentID,
		Type:         output.Type,
		Amount:       output.Amount,
		Balance:      output.Balance,
		Date:         output.Date,
	})
}

func mapMovementError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, appinvestment.ErrInvestmentNotFound):
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"message": err.Error()})
	case errors.Is(err, appinvestment.ErrInsufficientBalance):
		return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()})
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"message": "could not create movement",
		})
	}
}

func parseInvestmentID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"message": "invalid investment id",
		})
	}
	return uint(id), nil
}
