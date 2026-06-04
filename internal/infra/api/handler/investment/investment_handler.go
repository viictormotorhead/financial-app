package investment

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	apptag "github.com/viictormotorhead/financial-app/internal/application/tag"
	apihandler "github.com/viictormotorhead/financial-app/internal/infra/api/handler"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/request"
	investmentresponse "github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
	"github.com/viictormotorhead/financial-app/internal/infra/api/response"
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
)

type InvestmentHandlerIF interface {
	Create(c echo.Context) error
	List(c echo.Context) error
	GetByID(c echo.Context) error
	CreateMovement(c echo.Context) error
	CreateValuation(c echo.Context) error
}

type InvestmentHandler struct {
	createInvestment    use_cases.CreateInvestmentUseCaseIF
	listInvestments     use_cases.ListInvestmentsUseCaseIF
	getInvestmentDetail use_cases.GetInvestmentDetailUseCaseIF
	createMovement      use_cases.CreateMovementUseCaseIF
	createValuation     use_cases.CreateValuationUseCaseIF
}

func NewInvestmentHandler(
	createInvestment use_cases.CreateInvestmentUseCaseIF,
	listInvestments use_cases.ListInvestmentsUseCaseIF,
	getInvestmentDetail use_cases.GetInvestmentDetailUseCaseIF,
	createMovement use_cases.CreateMovementUseCaseIF,
	createValuation use_cases.CreateValuationUseCaseIF,
) InvestmentHandlerIF {
	return &InvestmentHandler{
		createInvestment:    createInvestment,
		listInvestments:     listInvestments,
		getInvestmentDetail: getInvestmentDetail,
		createMovement:      createMovement,
		createValuation:     createValuation,
	}
}

func (h *InvestmentHandler) Create(c echo.Context) error {
	var req request.CreateInvestmentRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(http.StatusBadRequest, "request body must be valid JSON")
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
		if httpErr := apihandler.HTTPErrorFromUseCase(err); httpErr != nil {
			return httpErr
		}
		var tagsNotFound *apptag.TagsNotFoundError
		if errors.As(err, &tagsNotFound) {
			return response.Error(http.StatusUnprocessableEntity, tagsNotFound.Error())
		}
		return response.Error(http.StatusInternalServerError, "could not create investment")
	}

	tags := output.Tags
	if tags == nil {
		tags = []string{}
	}

	return response.OK(c, http.StatusCreated, "investment created successfully", map[string]investmentresponse.InvestmentResponse{
		"investment": {
			ID:             output.ID,
			Name:           output.Name,
			Balance:        output.Balance,
			InitialBalance: output.InitialBalance,
			Tags:           tags,
			CreatedAt:      output.CreatedAt,
		},
	})
}
