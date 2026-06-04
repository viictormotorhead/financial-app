package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	"github.com/viictormotorhead/financial-app/internal/infra/api/response"
)

func HTTPErrorFromUseCase(err error) *echo.HTTPError {
	if errors.Is(err, appauth.ErrUnauthorized) {
		return response.Error(http.StatusUnauthorized, appauth.ErrUnauthorized.Error())
	}
	return nil
}
