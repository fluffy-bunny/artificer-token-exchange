package main

import (
	"fmt"
	"net/http"

	tex_config "tex/internal/contracts/config"
	tex_middleware_container "tex/internal/middleware/container"
	"tex/internal/shared"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/core"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

var version = "Development"

func main() {
	startup := NewStartup()
	config := &tex_config.Config{}
	err := core.LoadConfig(startup.GetConfigOptions())
	if err != nil {
		panic(err)
	}
	builder, _ := di.NewBuilder(di.App, di.Request, "transient")
	di.AddSingletonTypeByObj(builder, config)
	err = startup.ConfigureServices(builder)
	if err != nil {
		panic(err)
	}
	shared.RootContainer = builder.Build()

	e := echo.New()
	e.Use(tex_middleware_container.EnsureScopedContainer(shared.RootContainer))
	startup.Configure(e, shared.RootContainer)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
