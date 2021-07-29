package httputils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

// WrapBindError wraps given validation error to echo.HTTPError.
// If possible to cast validator.ValidationErrors, then use first error elts with http.StatusUnprocessableEntity
// otherwise return err's message with http.StatusInternalServerError
func WrapBindError(err error) error {
	switch v := err.(type) {
	case validator.ValidationErrors:
		msg := fmt.Sprintf("%s validation error. reason: %s", v[0].Field(), v[0].Tag())
		return NewError(http.StatusUnprocessableEntity, msg)
	case *echo.HTTPError:
		return v
	default:
		return NewError(http.StatusInternalServerError, err.Error())
	}
}

// NewUnauthorized returns echo.HTTPError with 401 Unauthorized and given message.
func NewUnauthorized() error {
	return NewError(http.StatusUnauthorized, "auth required")
}

// NewStatusUnprocessableEntity returns echo.HTTPError with 422 Unprocessable Entity and given message.
func NewStatusUnprocessableEntity(msg string) error {
	return NewError(http.StatusUnprocessableEntity, msg)
}

// NewNotFoundError returns echo.HTTPError with 404 Not found and given message.
// If empty message, then use "resource not found" as default message.
func NewNotFoundError(msg string) error {
	if msg == "" {
		msg = "resource not found"
	}
	return NewError(http.StatusNotFound, msg)
}

// NewInternalServerError returns echo.HTTPError with 500 internal server error and given error's message.
// If provide nil error, then use "interval server error occur" as default message.
func NewInternalServerError(err error) error {
	msg := "interval server error occur"
	if err != nil {
		msg = err.Error()
	}
	return NewError(http.StatusInternalServerError, msg)
}

func NewError(statusCode int, msg string) error {
	return echo.NewHTTPError(statusCode, &Error{
		Errors: map[string]interface{}{
			"body": msg,
		},
	})
}
