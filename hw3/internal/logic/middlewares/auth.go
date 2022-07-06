package middlewares

import (
	"crypto/subtle"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/simonnik/GB_Backend2_GO/hw3/internal/config"
)

type JwtCustomClaims struct {
	jwt.StandardClaims
}

func (c *JwtCustomClaims) VerifySubject(cmp string, req bool) bool {
	if c.Subject == "" {
		return !req
	}

	return subtle.ConstantTimeCompare([]byte(c.Subject), []byte(cmp)) != 0
}

func JWTAuthMiddleware(cfg *config.Config) []echo.MiddlewareFunc {
	var m []echo.MiddlewareFunc
	m = append(
		m,
		middleware.JWTWithConfig(middleware.JWTConfig{
			SuccessHandler: nil,
			SigningKey:     []byte(cfg.JWT.Secret),
			TokenLookup:    "header:X-API-KEY",
			Claims: &JwtCustomClaims{
				jwt.StandardClaims{},
			},
		}),
		jwtValidMiddleware(cfg),
	)

	return m
}

func jwtValidMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*JwtCustomClaims)
			vErr := new(jwt.ValidationError)

			if !claims.VerifyIssuer(cfg.JWT.Issuer, false) {
				vErr.Inner = fmt.Errorf("token used invalid issuer")
				vErr.Errors |= jwt.ValidationErrorIssuer
			}
			if !claims.VerifySubject(cfg.JWT.Subject, false) {
				vErr.Inner = fmt.Errorf("token used invalid subject")
				vErr.Errors |= jwt.ValidationErrorClaimsInvalid
			}

			// backward compatible error codes
			if vErr.Errors > 0 {
				return &echo.HTTPError{
					Code:     middleware.ErrJWTInvalid.Code,
					Message:  middleware.ErrJWTInvalid.Message,
					Internal: vErr,
				}
			}
			return next(c)
		}
	}
}
