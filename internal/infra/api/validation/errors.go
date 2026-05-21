package validation

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
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
		return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{
			Message: "invalid request",
			Errors:  []FieldError{{Field: "", Message: err.Error()}},
		})
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

	return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{
		Message: "validation failed",
		Errors:  fields,
	})
}
