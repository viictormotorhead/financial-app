package tag

import (
	"net/http"

	"github.com/labstack/echo/v4"

	tagresponse "github.com/viictormotorhead/financial-app/internal/infra/api/handler/tag/dto/response"
	"github.com/viictormotorhead/financial-app/internal/infra/api/response"
)

func (h *TagHandler) List(c echo.Context) error {
	output, err := h.listTags.List(c.Request().Context())
	if err != nil {
		return response.Error(http.StatusInternalServerError, "could not list tags")
	}

	tags := make([]tagresponse.TagResponse, 0, len(output.Tags))
	for _, tag := range output.Tags {
		tags = append(tags, tagresponse.TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
		})
	}

	return response.OK(c, http.StatusOK, "tags retrieved successfully", tags)
}
