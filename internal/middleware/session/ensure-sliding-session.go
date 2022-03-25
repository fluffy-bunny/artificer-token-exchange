package session

import (
	echostarter_wellknown "echo-starter/internal/wellknown"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"

	"echo-starter/internal/session"
)

func EnsureSlidingSession(root di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.GetSession(c)
			sess.Save(c.Request(), c.Response())
			subContainer, err := root.SubContainer()
			if err != nil {
				panic(err)
			}
			c.Set(echostarter_wellknown.SCOPED_CONTAINER_KEY, subContainer)
			return next(c)

		}
	}
}
