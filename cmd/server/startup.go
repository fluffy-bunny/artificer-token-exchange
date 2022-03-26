package main

import (
	echostarter_auth "echo-starter/internal/auth"
	tex_config "echo-starter/internal/contracts/config"

	services_auth_authenticator "echo-starter/internal/services/auth/authenticator"
	services_handlers_auth_callback "echo-starter/internal/services/handlers/auth/callback"
	services_handlers_auth_login "echo-starter/internal/services/handlers/auth/login"
	services_handlers_auth_logout "echo-starter/internal/services/handlers/auth/logout"
	services_handlers_deep "echo-starter/internal/services/handlers/deep"
	services_handlers_home "echo-starter/internal/services/handlers/home"

	middleware_claimsprincipal "echo-starter/internal/middleware/claimsprincipal"
	services_handler "echo-starter/internal/services/handler"

	services_claimsprovider "echo-starter/internal/services/claimsprovider"

	core_contracts "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/core"
	services_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/claimsprincipal"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Startup struct {
	config *tex_config.Config
}

func NewStartup() *Startup {
	return &Startup{
		config: &tex_config.Config{},
	}
}
func (s *Startup) GetPort() int {
	return s.config.Port
}
func (s *Startup) GetConfigOptions() *core_contracts.ConfigOptions {
	return &core_contracts.ConfigOptions{
		RootConfig:  []byte(tex_config.ConfigDefaultJSON),
		Destination: s.config,
	}
}
func (s *Startup) ConfigureServices(builder *di.Builder) error {
	di.AddSingletonTypeByObj(builder, s.config)
	services_auth_authenticator.AddSingletonIOIDCAuthenticator(builder)
	services_handlers_auth_login.AddScopedIHandler(builder)
	services_handlers_auth_callback.AddScopedIHandler(builder)
	services_handlers_auth_logout.AddScopedIHandler(builder)
	services_handlers_home.AddScopedIHandler(builder)
	services_handlers_deep.AddScopedIHandler(builder)
	services_handler.AddSingletonIHandlerFactory(builder)
	services_core_claimsprincipal.AddScopedIClaimsPrincipal(builder)
	services_claimsprovider.AddSingletonIClaimsProviderMock(builder)
	return nil
}
func (s *Startup) Configure(e *echo.Echo, root di.Container) error {
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			id := uuid.New()
			return id.String()
		},
	}))
	// DevelopmentMiddlewareUsingClaimsMap adds all the needed claims so that FinalAuthVerificationMiddlewareUsingClaimsMap succeeds
	//e.Use(middleware_claimsprincipal.DevelopmentMiddlewareUsingClaimsMap(echostarter_auth.BuildGrpcEntrypointPermissionsClaimsMap(), true))

	e.Use(middleware_claimsprincipal.AuthenticatedSessionToClaimsPrincipalMiddleware())
	e.Use(middleware_claimsprincipal.FinalAuthVerificationMiddlewareUsingClaimsMap(echostarter_auth.BuildGrpcEntrypointPermissionsClaimsMap(), true))
	return nil
}
