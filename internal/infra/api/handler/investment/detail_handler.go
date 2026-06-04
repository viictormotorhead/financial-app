package investment

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	appinvestment "github.com/viictormotorhead/financial-app/internal/application/investment"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	apihandler "github.com/viictormotorhead/financial-app/internal/infra/api/handler"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
	apiresponse "github.com/viictormotorhead/financial-app/internal/infra/api/response"
)

func (h *InvestmentHandler) GetByID(c echo.Context) error {
	investmentID, err := parseInvestmentID(c)
	if err != nil {
		return err
	}

	output, err := h.getInvestmentDetail.Get(c.Request().Context(), use_cases.GetInvestmentDetailQuery{
		InvestmentID: investmentID,
	})
	if err != nil {
		if httpErr := apihandler.HTTPErrorFromUseCase(err); httpErr != nil {
			return httpErr
		}
		if errors.Is(err, appinvestment.ErrInvestmentNotFound) {
			return apiresponse.Error(http.StatusNotFound, err.Error())
		}
		return apiresponse.Error(http.StatusInternalServerError, "could not get investment detail")
	}

	tags := output.Tags
	if tags == nil {
		tags = []string{}
	}

	series := make([]response.InvestmentSeriesPointResponse, 0, len(output.Series))
	for _, point := range output.Series {
		series = append(series, response.InvestmentSeriesPointResponse{
			Date:  formatDate(point.Date),
			Value: point.Value,
		})
	}

	movements := make([]response.InvestmentDetailMovementResponse, 0, len(output.Movements))
	for _, movement := range output.Movements {
		movements = append(movements, response.InvestmentDetailMovementResponse{
			ID:           movement.ID,
			Date:         formatDate(movement.Date),
			Type:         movement.Type,
			Amount:       movement.Amount,
			BalanceAfter: movement.BalanceAfter,
		})
	}

	return apiresponse.OK(c, http.StatusOK, "", response.InvestmentDetailResponse{
		ID:            output.ID,
		Name:          output.Name,
		Tags:          tags,
		CurrentValue:  output.CurrentValue,
		GrowthPercent: output.GrowthPercent,
		GrowthAmount:  output.GrowthAmount,
		Series:        series,
		Movements:     movements,
	})
}

func formatDate(value time.Time) string {
	return value.UTC().Format("2006-01-02")
}
