package session

import (
	contracts_auth "echo-starter/internal/contracts/auth"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

func EnsureSlidingAuthCookie(root di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// get authCookie service once during configuration
		authCookie := contracts_auth.GetIAuthCookieFromContainer(root)
		return func(c echo.Context) error {
			authCookie.RefreshAuthCookie(c)
			return next(c)
		}
	}
}
