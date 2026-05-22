package validation

import (
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"

	"github.com/viictormotorhead/financial-app/internal/infra/api/response"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

func HTTPError(err error) *echo.HTTPError {
	errs, ok := err.(validation.Errors)
	if !ok {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	fields := make([]FieldError, 0, len(errs))
	for field, fieldErr := range errs {
		if fieldErr == nil {
			continue
		}
		fields = append(fields, FieldError{
			Field:   field,
			Message: fieldErr.Error(),
		})
	}

	return response.Error(http.StatusBadRequest, formatMessage(ErrorResponse{
		Message: "validation failed",
		Errors:  fields,
	}))
}

func formatMessage(resp ErrorResponse) string {
	if len(resp.Errors) == 0 {
		if resp.Message != "" {
			return resp.Message
		}
		return "validation failed"
	}

	parts := make([]string, 0, len(resp.Errors))
	for _, fieldErr := range resp.Errors {
		if fieldErr.Field == "" {
			parts = append(parts, fieldErr.Message)
			continue
		}
		parts = append(parts, fmt.Sprintf("%s: %s", fieldErr.Field, fieldErr.Message))
	}

	if resp.Message != "" {
		return fmt.Sprintf("%s (%s)", resp.Message, strings.Join(parts, "; "))
	}

	return strings.Join(parts, "; ")
}
