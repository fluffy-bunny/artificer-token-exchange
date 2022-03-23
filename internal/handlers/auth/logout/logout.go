package logout

import (
	"net/http"

	auth_authenticator "echo-starter/internal/auth/authenticator"
	"echo-starter/internal/session"

	"github.com/labstack/echo/v4"
)

func Handler(auth *auth_authenticator.Authenticator) func(c echo.Context) error {
	return func(c echo.Context) error {
		session.TerminateSession(c)

		// Redirect to home page.
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return nil
	}
}
