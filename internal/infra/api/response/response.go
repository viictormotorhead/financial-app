package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	StatusOK    = "ok"
	StatusError = "error"
)

type Envelope struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func OK(c echo.Context, code int, message string, data any) error {
	return c.JSON(code, Envelope{
		Status:  StatusOK,
		Message: message,
		Data:    data,
	})
}

func Error(code int, message string) *echo.HTTPError {
	return echo.NewHTTPError(code, Envelope{
		Status:  StatusError,
		Message: message,
		Data:    nil,
	})
}

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := http.StatusText(code)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = messageFromBody(he.Message, code)
	}

	_ = c.JSON(code, Envelope{
		Status:  StatusError,
		Message: message,
		Data:    nil,
	})
}

func messageFromBody(body any, code int) string {
	switch v := body.(type) {
	case Envelope:
		if v.Message != "" {
			return v.Message
		}
	case *Envelope:
		if v != nil && v.Message != "" {
			return v.Message
		}
	case string:
		if v != "" {
			return v
		}
	case map[string]string:
		if msg := v["message"]; msg != "" {
			return msg
		}
	case map[string]any:
		if msg, ok := v["message"].(string); ok && msg != "" {
			return msg
		}
	}

	if msg := http.StatusText(code); msg != "" {
		return msg
	}

	return "unexpected error"
}
