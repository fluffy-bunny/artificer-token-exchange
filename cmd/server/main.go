package main

import (
	"fmt"
	"net/http"

	tex_middleware_container "echo-starter/internal/middleware/container"
	"echo-starter/internal/shared"
	"os"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/core"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var version = "Development"

func main() {
	d := os.Getenv("OIDC__CLIENT_ID")
	fmt.Println(d)

	d = os.Getenv("AUTH0_CLIENT_ID")
	fmt.Println(d)
	startup := NewStartup()
	configOptions := startup.GetConfigOptions()
	err := core.LoadConfig(configOptions)
	if err != nil {
		panic(err)
	}
	builder, _ := di.NewBuilder(di.App, di.Request, "transient")
	err = startup.ConfigureServices(builder)
	if err != nil {
		panic(err)
	}
	shared.RootContainer = builder.Build()

	e := echo.New()
	e.Use(tex_middleware_container.EnsureScopedContainer(shared.RootContainer))
	e.Use(middleware.Logger())
	startup.Configure(e, shared.RootContainer)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	port := startup.GetPort()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}
