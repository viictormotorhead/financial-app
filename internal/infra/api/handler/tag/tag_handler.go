package tag

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	apptag "github.com/viictormotorhead/financial-app/internal/application/tag"
	"github.com/viictormotorhead/financial-app/internal/application/tag/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/tag/dto/request"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/tag/dto/response"
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
)

type TagHandlerIF interface {
	Create(c echo.Context) error
}

type TagHandler struct {
	createTag use_cases.CreateTagUseCaseIF
}

func NewTagHandler(createTag use_cases.CreateTagUseCaseIF) TagHandlerIF {
	return &TagHandler{createTag: createTag}
}

func (h *TagHandler) Create(c echo.Context) error {
	var req request.CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, apivalidation.ErrorResponse{
			Message: "invalid request body",
			Errors:  []apivalidation.FieldError{{Field: "", Message: "request body must be valid JSON"}},
		})
	}

	if err := req.Validate(); err != nil {
		return apivalidation.HTTPError(err)
	}

	output, err := h.createTag.Create(c.Request().Context(), use_cases.CreateTagCommand{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, apptag.ErrTagNameAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, map[string]string{
				"message": err.Error(),
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"message": "could not create tag",
		})
	}

	return c.JSON(http.StatusCreated, response.TagResponse{
		ID:          output.ID,
		Name:        output.Name,
		Description: output.Description,
	})
}
