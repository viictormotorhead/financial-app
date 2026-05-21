package investment

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	apptag "github.com/viictormotorhead/financial-app/internal/application/tag"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/request"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
)

type InvestmentHandlerIF interface {
	Create(c echo.Context) error
	CreateMovement(c echo.Context) error
	CreateValuation(c echo.Context) error
}

type InvestmentHandler struct {
	createInvestment use_cases.CreateInvestmentUseCaseIF
	createMovement   use_cases.CreateMovementUseCaseIF
	createValuation  use_cases.CreateValuationUseCaseIF
}

func NewInvestmentHandler(
	createInvestment use_cases.CreateInvestmentUseCaseIF,
	createMovement use_cases.CreateMovementUseCaseIF,
	createValuation use_cases.CreateValuationUseCaseIF,
) InvestmentHandlerIF {
	return &InvestmentHandler{
		createInvestment: createInvestment,
		createMovement:   createMovement,
		createValuation:  createValuation,
	}
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
		var tagsNotFound *apptag.TagsNotFoundError
		if errors.As(err, &tagsNotFound) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]interface{}{
				"message": tagsNotFound.Error(),
				"tags":    tagsNotFound.Missing,
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"message": "could not create investment",
		})
	}

	tags := output.Tags
	if tags == nil {
		tags = []string{}
	}

	return c.JSON(http.StatusCreated, response.InvestmentResponse{
		ID:             output.ID,
		Name:           output.Name,
		Balance:        output.Balance,
		InitialBalance: output.InitialBalance,
		Tags:           tags,
		CreatedAt:      output.CreatedAt,
	})
}
