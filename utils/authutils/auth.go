package authutils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	AuthScheme = "Token"
)

type JWTClaims struct {
	UserID uint
	jwt.StandardClaims
}

func MakeJWTToken(userID uint, secret []byte, expires time.Duration) (string, error) {
	c := &JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expires).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

// CurrentUser returns current user id which stored at echo.Context if exist, otherwise returns 0.
func CurrentUser(ctx echo.Context) uint {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return 0
	}
	return token.Claims.(*JWTClaims).UserID
}

func SetAuthToken(r *http.Request, token string) {
	r.Header.Set("Authorization", fmt.Sprintf("%s %s", AuthScheme, token))
}
