package group

import (
	"github.com/labstack/echo/v4"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/tag"
)

const TagPathV1 = "/v1/tags"

type TagRoutes struct{}

func NewTagRoutes(group *echo.Group, handler tag.TagHandlerIF) *TagRoutes {
	routes := group.Group(TagPathV1)
	routes.POST("/", handler.Create)

	return &TagRoutes{}
}
