package httputils

import "github.com/labstack/echo/v4"

// BindAndValidate bind and validate from given echo.Context and interface v.
func BindAndValidate(ctx echo.Context, v interface{}) error {
	if err := ctx.Bind(v); err != nil {
		return err
	}
	if err := ctx.Validate(v); err != nil {
		return err
	}
	return nil
}
