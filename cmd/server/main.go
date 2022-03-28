package main

import (
	"echo-starter/internal/templates"
	"fmt"

	contracts_config "echo-starter/internal/contracts/config"
	services_container "echo-starter/internal/services/container"

	contracts_handler "echo-starter/internal/contracts/handler"
	middleware_container "echo-starter/internal/middleware/container"
	middleware_session "echo-starter/internal/middleware/session"

	"echo-starter/internal/shared"
	echostarter_utils "echo-starter/internal/utils"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/core"
	"github.com/google/uuid"

	"echo-starter/internal/wellknown"

	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var version = "Development"

func main() {
	appInstanceID := uuid.New().String()

	fmt.Println("Version:" + version)
	DumpPath("./")
	var err error
	startup := NewStartup()
	configOptions := startup.GetConfigOptions()
	err = core.LoadConfig(configOptions)
	if err != nil {
		panic(err)
	}
	appConfig := configOptions.Destination.(*contracts_config.Config)
	if core_utils.IsEmptyOrNil(appConfig.Oidc.CallbackURL) {
		appConfig.Oidc.CallbackURL = fmt.Sprintf("http://localhost:%v%s",
			appConfig.Port,
			wellknown.OIDCCallbackPath)
	}
	fmt.Println(echostarter_utils.PrettyJSON(appConfig))
	builder, _ := di.NewBuilder(di.App, di.Request, "transient")
	services_container.AddContainerAccessorFunc(builder)
	err = startup.ConfigureServices(builder)
	if err != nil {
		panic(err)
	}
	shared.RootContainer = builder.Build()

	e := echo.New()
	//Set Renderer
	e.Renderer = templates.GetTemplateRender("./templates")

	if core_utils.IsEmptyOrNil(appConfig.SessionKey) {
		fmt.Println("WARNING: SESSION_KEY must be set for production......")
		appConfig.SessionKey = core_utils.RandomString(32)
		fmt.Printf("SESSION_KEY: %v\n", appConfig.SessionKey)
	}
	if core_utils.IsEmptyOrNil(appConfig.SessionEncryptionKey) {
		fmt.Println("WARNING: SESSION_ENCRYPTION_KEY must be set for production......")
		appConfig.SessionEncryptionKey = core_utils.RandomString(32)
		fmt.Printf("SESSION_ENCRYPTION_KEY: %v\n", appConfig.SessionEncryptionKey)
	}

	e.Use(middleware.Logger())
	e.Use(middleware_container.EnsureScopedContainer(shared.RootContainer))

	// we don't have a shared backend session store (i.e. redis), so fat cookies it is
	sessionStore := sessions.NewCookieStore([]byte(appConfig.SessionKey), []byte(appConfig.SessionEncryptionKey))
	sessionStore.Options.MaxAge = appConfig.SessionMaxAgeSeconds

	e.Use(session.Middleware(sessionStore))
	e.Use(middleware_session.EnsureSlidingSession(shared.RootContainer))
	if appConfig.ApplicationEnvironment == contracts_config.Environment_Development {
		e.Use(middleware_session.EnsureDevelopmentSession(appInstanceID))
	}
	e.Use(middleware_session.EnsureSlidingAuthCookie(shared.RootContainer))

	app := e.Group("")
	app.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))
	startup.Configure(e, shared.RootContainer)
	e.Use(middleware.Recover())

	app.Static("/css", "./css")
	app.Static("/assets", "./assets")
	app.Static("/js", "./js")
	handlerFactory := contracts_handler.GetIHandlerFactoryFromContainer(shared.RootContainer)
	handlerFactory.RegisterHandlers(app)

	port := startup.GetPort()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}
