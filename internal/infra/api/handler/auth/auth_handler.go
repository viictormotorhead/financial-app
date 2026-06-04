package auth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	appuser "github.com/viictormotorhead/financial-app/internal/application/user"
	"github.com/viictormotorhead/financial-app/internal/application/user/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/user/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/auth/dto/request"
	authresponse "github.com/viictormotorhead/financial-app/internal/infra/api/handler/auth/dto/response"
	"github.com/viictormotorhead/financial-app/internal/infra/api/response"
	apivalidation "github.com/viictormotorhead/financial-app/internal/infra/api/validation"
)

type AuthHandlerIF interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

type AuthHandler struct {
	register use_cases.RegisterUserUseCaseIF
	login    use_cases.LoginUserUseCaseIF
}

func NewAuthHandler(
	register use_cases.RegisterUserUseCaseIF,
	login use_cases.LoginUserUseCaseIF,
) AuthHandlerIF {
	return &AuthHandler{
		register: register,
		login:    login,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req request.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(http.StatusBadRequest, "request body must be valid JSON")
	}

	if err := req.Validate(); err != nil {
		return apivalidation.HTTPError(err)
	}

	output, err := h.register.Register(c.Request().Context(), use_cases.RegisterUserCommand{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, appuser.ErrUsernameAlreadyExists) {
			return response.Error(http.StatusConflict, err.Error())
		}
		return response.Error(http.StatusInternalServerError, "could not register user")
	}

	return response.OK(c, http.StatusCreated, "user registered successfully", map[string]authresponse.AuthResponse{
		"auth": toAuthResponse(output),
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(http.StatusBadRequest, "request body must be valid JSON")
	}

	if err := req.Validate(); err != nil {
		return apivalidation.HTTPError(err)
	}

	output, err := h.login.Login(c.Request().Context(), use_cases.LoginUserCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, appauth.ErrInvalidCredentials) {
			return response.Error(http.StatusUnauthorized, err.Error())
		}
		return response.Error(http.StatusInternalServerError, "could not login")
	}

	return response.OK(c, http.StatusOK, "login successful", map[string]authresponse.AuthResponse{
		"auth": toAuthResponse(output),
	})
}

func toAuthResponse(output outputs.AuthOutputDTO) authresponse.AuthResponse {
	return authresponse.AuthResponse{
		Token: output.Token,
		User: authresponse.AuthUserResponse{
			UserID:   output.User.UserID,
			Name:     output.User.Name,
			Username: output.User.Username,
		},
	}
}
