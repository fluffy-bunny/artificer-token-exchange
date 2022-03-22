package container

import (
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

func EnsureScopedContainer(root di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			subContainer, err := root.SubContainer()
			if err != nil {
				panic(err)
			}
			c.Set("container", subContainer)
			return next(c)
		}
	}
}
