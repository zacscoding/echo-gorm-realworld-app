package httputils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

// WrapBindError wraps given error to echo.HTTPError
func WrapBindError(err error) error {
	// TODO: implement me
	return err
}

// NewUnauthorized returns echo.HTTPError with 401 Unauthorized and given message.
func NewUnauthorized() error {
	return newError(http.StatusUnauthorized, "auth required")
}

// NewStatusUnprocessableEntity returns echo.HTTPError with 422 Unprocessable Entity and given message.
func NewStatusUnprocessableEntity(msg string) error {
	return newError(http.StatusUnprocessableEntity, msg)
}

// NewNotFoundError returns echo.HTTPError with 404 Not found and given message.
// If empty message, then use "resource not found" as default message.
func NewNotFoundError(msg string) error {
	if msg == "" {
		msg = "resource not found"
	}
	return newError(http.StatusNotFound, msg)
}

// NewInternalServerError returns echo.HTTPError with 500 internal server error and given error's message.
// If provide nil error, then use "interval server error occur" as default message.
func NewInternalServerError(err error) error {
	msg := "interval server error occur"
	if err != nil {
		msg = err.Error()
	}
	return newError(http.StatusInternalServerError, msg)
}

func newError(statusCode int, msg string) error {
	return echo.NewHTTPError(statusCode, &Error{
		Errors: map[string]interface{}{
			"body": msg,
		},
	})
}
