package startup

import (
	"net"

	core_contracts "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/core"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	Options struct {
		Listener net.Listener
	}

	IStartup interface {
		GetOptions() *Options
		GetConfigOptions() *core_contracts.ConfigOptions
		GetPort() int
		ConfigureServices(builder *di.Builder) error
		Configure(e *echo.Echo, root di.Container) error
	}
)
