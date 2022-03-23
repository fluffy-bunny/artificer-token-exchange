package main

import (
	"fmt"
	"io"
	"text/template"

	auth_authenticator "echo-starter/internal/auth/authenticator"
	contracts_config "echo-starter/internal/contracts/config"
	handlers_auth_callback "echo-starter/internal/handlers/auth/callback"
	handlers_auth_login "echo-starter/internal/handlers/auth/login"
	handlers_auth_logout "echo-starter/internal/handlers/auth/logout"
	handlers_home "echo-starter/internal/handlers/home"

	handlers_auth_user "echo-starter/internal/handlers/auth/user"

	middleware_container "echo-starter/internal/middleware/container"

	"echo-starter/internal/shared"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/core"
	"github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var version = "Development"

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
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
	//Set Renderer
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.tpl")),
	}
	e.Renderer = renderer

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
	e.Use(middleware.Recover())

	app := e.Group("")
	app.Static("/css", "./assets/css")
	app.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))
	// we don't have a shared backend session store (i.e. redis), so fat cookies it is
	app.Use(session.Middleware(sessions.NewCookieStore([]byte(appConfig.SessionKey), []byte(appConfig.SessionEncryptionKey))))

	startup.Configure(e, shared.RootContainer)

	e.GET("/", handlers_home.Handler())
	authenticator, err := auth_authenticator.New(appConfig)
	if err != nil {
		panic(err)
	}
	app.GET("/oidc", handlers_auth_callback.Handler(authenticator))
	app.GET("/login", handlers_auth_login.Handler(authenticator))
	app.GET("/user", handlers_auth_user.Handler())
	app.GET("/logout", handlers_auth_logout.Handler(authenticator))
	port := startup.GetPort()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}
