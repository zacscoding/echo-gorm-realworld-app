package authutils

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/utils/httputils"
)

// NewJWTMiddleware returns JWT auth middleware with given optional paths and secret.
// requests in optionalAuthPaths will skip middleware if header's Authorization is empty.
func NewJWTMiddleware(optionalAuthPaths map[string]struct{}, secret string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(
		middleware.JWTConfig{
			Skipper: func(ctx echo.Context) bool {
				if _, ok := optionalAuthPaths[ctx.Path()]; !ok {
					return false
				}
				return ctx.Request().Header.Get("Authorization") == ""
			},
			TokenLookup: "header:Authorization",
			Claims:      &JWTClaims{},
			SigningKey:  []byte(secret),
			AuthScheme:  AuthScheme,
			ErrorHandlerWithContext: func(err error, ctx echo.Context) error {
				logging.FromContext(ctx.Request().Context()).Errorw("auth failed", "err", err)
				return httputils.NewUnauthorized()
			},
		},
	)
}
