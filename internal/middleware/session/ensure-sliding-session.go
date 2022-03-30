package session

import (
	"echo-starter/internal/session"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

func EnsureSlidingSession(container di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			sess := session.GetSession(c)
			if !sess.IsNew {
				// we don't want to create a new session if nobody every created one before
				// we are only here to ensure that the session is an old one and slide it out.
				// i.e. bump out the expiration time
				sess.Save(c.Request(), c.Response())
			}
			authSession := session.GetAuthSession(c)
			if !authSession.IsNew {
				// we don't want to create a new session if nobody every created one before
				// we are only here to ensure that the session is an old one and slide it out.
				// i.e. bump out the expiration time
				authSession.Save(c.Request(), c.Response())
			}
			return next(c)
		}
	}
}
