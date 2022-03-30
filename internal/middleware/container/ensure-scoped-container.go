package container

import (
	contracts_contextaccessor "echo-starter/internal/contracts/contextaccessor"
	echostarter_wellknown "echo-starter/internal/wellknown"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"

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
			c.Set(echostarter_wellknown.SCOPED_CONTAINER_KEY, subContainer)
			internalContextAccessor := contracts_contextaccessor.GetIInternalEchoContextAccessorFromContainer(subContainer)
			internalContextAccessor.SetContext(c)
			contracts_logger.GetILoggerFromContainer(subContainer) // force an instantiation
			/*
				contextAccessor := contracts_contextaccessor.GetIEchoContextAccessorFromContainer(subContainer)
				d := contextAccessor.GetContext()
				if d != c {
					panic("contextAccessor.GetContext() != c")
				}
			*/
			return next(c)
		}
	}
}
