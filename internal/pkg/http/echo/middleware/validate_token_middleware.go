package middleware

import (
	"fmt"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// ValidateBearerToken from request
func ValidateBearerToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Parse and verify jwt access token
			auth, ok := bearerAuth(c.Request())
			if !ok {
				return fmt.Errorf("parse jwt access token error")
			}
			token, err := jwt.ParseWithClaims(auth, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("parse signing method error")
				}
				return []byte("secret"), nil
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			c.Set("token", token)
			return next(c)
		}
	}
}

// BearerAuth parse bearer token
func bearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access_token")
	}
	return token, token != ""
}
