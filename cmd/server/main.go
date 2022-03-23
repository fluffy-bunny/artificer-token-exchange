package main

import (
	"fmt"
	"net/http"

	auth_authenticator "echo-starter/internal/auth/authenticator"
	contracts_config "echo-starter/internal/contracts/config"
	handlers_auth_callback "echo-starter/internal/handlers/auth/callback"
	handlers_auth_login "echo-starter/internal/handlers/auth/login"
	handlers_auth_logout "echo-starter/internal/handlers/auth/logout"

	handlers_auth_user "echo-starter/internal/handlers/auth/user"

	middleware_container "echo-starter/internal/middleware/container"

	"echo-starter/internal/shared"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/core"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var version = "Development"

func main() {
	fmt.Println("Version:" + version)
	startup := NewStartup()
	configOptions := startup.GetConfigOptions()
	err := core.LoadConfig(configOptions)
	if err != nil {
		panic(err)
	}
	appConfig := configOptions.Destination.(*contracts_config.Config)
	if core_utils.IsEmptyOrNil(appConfig.Oidc.CallbackURL) {
		appConfig.Oidc.CallbackURL = fmt.Sprintf("http://localhost:%v/oidc", appConfig.Port)
	}
	builder, _ := di.NewBuilder(di.App, di.Request, "transient")
	err = startup.ConfigureServices(builder)
	if err != nil {
		panic(err)
	}
	shared.RootContainer = builder.Build()

	e := echo.New()

	if core_utils.IsEmptyOrNil(appConfig.SessionKey) {
		fmt.Println("WARNING: SESSION_KEY must be set for production......")
		appConfig.SessionKey = utils.RandomString(32)
		fmt.Printf("SESSION_KEY: %v\n", appConfig.SessionKey)
	}
	if core_utils.IsEmptyOrNil(appConfig.SessionEncryptionKey) {
		fmt.Println("WARNING: SESSION_ENCRYPTION_KEY must be set for production......")
		appConfig.SessionEncryptionKey = utils.RandomString(32)
		fmt.Printf("SESSION_ENCRYPTION_KEY: %v\n", appConfig.SessionEncryptionKey)
	}
	// we don't have a shared backend session store (i.e. redis), so fat cookies it is
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(appConfig.SessionKey), []byte(appConfig.SessionEncryptionKey))))

	e.Use(middleware_container.EnsureScopedContainer(shared.RootContainer))
	e.Use(middleware.Logger())
	startup.Configure(e, shared.RootContainer)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	authenticator, err := auth_authenticator.New(appConfig)
	if err != nil {
		panic(err)
	}
	e.GET("/oidc", handlers_auth_callback.Handler(authenticator))
	e.GET("/login", handlers_auth_login.Handler(authenticator))
	e.GET("/user", handlers_auth_user.Handler(authenticator))
	e.GET("/logout", handlers_auth_logout.Handler(authenticator))
	port := startup.GetPort()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}
