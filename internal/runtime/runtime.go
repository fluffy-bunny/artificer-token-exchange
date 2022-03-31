package runtime

import (
	contracts_config "echo-starter/internal/contracts/config"
	contracts_container "echo-starter/internal/contracts/container"
	contracts_handler "echo-starter/internal/contracts/handler"
	contracts_startup "echo-starter/internal/contracts/startup"
	"echo-starter/internal/templates"
	"echo-starter/internal/wellknown"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"

	services_timeutils "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/timeutils"

	middleware_container "echo-starter/internal/middleware/container"
	middleware_logger "echo-starter/internal/middleware/logger"

	middleware_session "echo-starter/internal/middleware/session"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/core"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/ziflex/lecho"
)

type (
	Runtime struct {
		Startup    contracts_startup.IStartup
		Config     *contracts_config.Config
		Container  di.Container
		e          *echo.Echo
		instanceID string
	}
)

func New(startup contracts_startup.IStartup) *Runtime {
	return &Runtime{
		Startup:    startup,
		instanceID: uuid.New().String(),
	}
}
func (s *Runtime) GetContainer() di.Container {
	return s.Container
}
func (s *Runtime) Run() error {
	startupOptions := s.Startup.GetOptions()
	configOptions := s.Startup.GetConfigOptions()
	err := core.LoadConfig(configOptions)
	if err != nil {
		return err
	}
	s.Config = configOptions.Destination.(*contracts_config.Config)
	if core_utils.IsEmptyOrNil(s.Config.Oidc.CallbackURL) {
		s.Config.Oidc.CallbackURL = fmt.Sprintf("http://localhost:%v%s",
			s.Config.Port,
			wellknown.OIDCCallbackPath)
	}
	if s.Config.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch strings.ToLower(s.Config.LogLevel) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	builder, _ := di.NewBuilder(di.App, di.Request, "transient")
	err = s.AddDefaultServices(builder)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add default services")
	}
	err = s.Startup.ConfigureServices(builder)
	if err != nil {
		log.Error().Err(err).Msg("Failed to configure services")
	}
	s.Container = builder.Build()
	s.Startup.SetContainer(s.Container)
	s.e = echo.New()
	//use our own zerolog logger
	s.e.Logger = lecho.New(os.Stdout)
	//Set Renderer
	s.e.Renderer = templates.GetTemplateRender("./templates")

	// SECURE COOKIE
	if core_utils.IsEmptyOrNil(s.Config.SecureCookieHashKey) {
		fmt.Println("WARNING: SECURE_COOKIE_HASH_KEY must be set for production......")
		key := securecookie.GenerateRandomKey(32)
		encodedString := base64.StdEncoding.EncodeToString(key)
		s.Config.SecureCookieHashKey = encodedString
		fmt.Printf("SECURE_COOKIE_HASH_KEY: %v\n", s.Config.SecureCookieHashKey)
	}
	if core_utils.IsEmptyOrNil(s.Config.SecureCookieEncryptionKey) {
		fmt.Println("WARNING: SECURE_COOKIE_ENCRYPTION_KEY must be set for production......")
		key := securecookie.GenerateRandomKey(32)
		encodedString := base64.StdEncoding.EncodeToString(key)
		s.Config.SecureCookieEncryptionKey = encodedString
		fmt.Printf("SECURE_COOKIE_ENCRYPTION_KEY: %v\n", s.Config.SecureCookieEncryptionKey)
	}
	// MIDDELWARE
	//-------------------------------------------------------
	s.e.Use(middleware_logger.EnsureContextLogger(s.Container))
	s.e.Use(middleware_logger.EnsureContextLoggerCorrelation(s.Container))
	s.e.Use(middleware_container.EnsureScopedContainer(s.Container))
	s.e.Use(session.Middleware(s.Startup.GetSessionStore()))
	s.e.Use(middleware_session.EnsureSlidingSession(s.Container))
	if s.Config.ApplicationEnvironment == contracts_config.Environment_Development {
		// this wipes out the session if we have a mismatch
		s.e.Use(middleware_session.EnsureDevelopmentSession(s.instanceID))
	}
	app := s.e.Group("")
	app.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-Csrf-Token,form:csrf",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))

	// we have all our required upfront middleware running
	// now we can add the optional startup ones.
	s.Startup.Configure(s.e, s.Container)

	// our middleware that runs at the end
	//-------------------------------------------------------
	s.e.Use(middleware.Recover())
	s.Startup.RegisterStaticRoutes(s.e)

	// register our handlers
	handlerFactory := contracts_handler.GetIHandlerFactoryFromContainer(s.Container)
	handlerFactory.RegisterHandlers(app)
	handlerDefinitions := contracts_handler.GetIHandlerDefinitions(s.Container)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Verbs", "Path"})

	for _, handlerDefinition := range handlerDefinitions {
		metaData := handlerDefinition.MetaData
		httpVerbs, _ := metaData["httpVerbs"].([]contracts_handler.HTTPVERB)
		verbBldr := strings.Builder{}

		for idx, verb := range httpVerbs {
			verbBldr.WriteString(verb.String())
			if idx < len(httpVerbs)-1 {
				verbBldr.WriteString(",")
			}
		}
		path, _ := metaData["path"].(string)

		t.AppendRow([]interface{}{verbBldr.String(), string(path)})

	}
	t.Render()
	// Finally start the server
	//----------------------------------------------------------------------------------
	port := s.Startup.GetPort()
	address := fmt.Sprintf(":%v", port)
	if startupOptions != nil && startupOptions.Listener != nil {
		// if we are here we are usually under test
		s.e.Listener = startupOptions.Listener
	}

	err = s.e.Start(address)
	if err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
	return err
}
func (s *Runtime) AddDefaultServices(builder *di.Builder) error {
	contracts_container.AddContainerAccessorFunc(builder, s.GetContainer)
	services_timeutils.AddTimeNow(builder)
	services_timeutils.AddTimeParse(builder)
	return nil
}
