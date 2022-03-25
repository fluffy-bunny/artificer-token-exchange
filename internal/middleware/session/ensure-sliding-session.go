package session

import (
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"

	"echo-starter/internal/session"
)

func EnsureSlidingSession(root di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.GetSession(c)
			sess.Save(c.Request(), c.Response())
			return next(c)
		}
	}
}
