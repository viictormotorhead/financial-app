package investment

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment/dto/response"
	apiresponse "github.com/viictormotorhead/financial-app/internal/infra/api/response"
)

func (h *InvestmentHandler) List(c echo.Context) error {
	output, err := h.listInvestments.List(c.Request().Context(), use_cases.ListInvestmentsQuery{
		Tags: parseTagsQueryParam(c),
	})
	if err != nil {
		return apiresponse.Error(http.StatusInternalServerError, "could not list investments")
	}

	items := make([]response.InvestmentAllocationResponse, 0, len(output.Items))
	for _, item := range output.Items {
		tags := item.Tags
		if tags == nil {
			tags = []string{}
		}

		items = append(items, response.InvestmentAllocationResponse{
			Investment:        item.Investment,
			Amount:            item.Amount,
			Percentage:        item.Percentage,
			PercentageGrowing: item.PercentageGrowing,
			Tags:              tags,
		})
	}

	return apiresponse.OK(c, http.StatusOK, "investments retrieved successfully", items)
}

func parseTagsQueryParam(c echo.Context) []string {
	rawTags := c.QueryParams()["tags"]
	if len(rawTags) == 0 {
		return nil
	}

	tags := make([]string, 0, len(rawTags))
	for _, value := range rawTags {
		for _, part := range strings.Split(value, ",") {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	return tags
}
