package main

import (
	echostarter_auth "echo-starter/internal/auth"
	tex_config "echo-starter/internal/contracts/config"
	services_auth_authenticator "echo-starter/internal/services/auth/authenticator"
	services_cookies "echo-starter/internal/services/cookies"
	services_handlers_about "echo-starter/internal/services/handlers/about"
	"echo-starter/internal/shared"

	core_services_timeutils "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/timeutils"

	// ACCOUNTS
	//----------------------------------------------------------------------------------------------------------------------
	services_handlers_accounts "echo-starter/internal/services/handlers/accounts"
	services_handlers_api_accounts "echo-starter/internal/services/handlers/api/accounts"

	// ARTISTS
	//----------------------------------------------------------------------------------------------------------------------
	services_handlers_api_artists "echo-starter/internal/services/handlers/api/artists"
	services_handlers_api_artists_artist "echo-starter/internal/services/handlers/api/artists/artist"
	services_handlers_api_artists_artist_albums "echo-starter/internal/services/handlers/api/artists/artist/albums"
	services_handlers_artists "echo-starter/internal/services/handlers/artists"

	services_handlers_auth_callback "echo-starter/internal/services/handlers/auth/callback"
	services_handlers_auth_login "echo-starter/internal/services/handlers/auth/login"
	services_handlers_auth_logout "echo-starter/internal/services/handlers/auth/logout"

	middleware_claimsprincipal "echo-starter/internal/middleware/claimsprincipal"
	middleware_session "echo-starter/internal/middleware/session"
	services_claimsprovider "echo-starter/internal/services/claimsprovider"
	services_handler "echo-starter/internal/services/handler"
	services_handlers_auth_profiles "echo-starter/internal/services/handlers/auth/profiles"
	services_handlers_auth_unauthorized "echo-starter/internal/services/handlers/auth/unauthorized"
	services_handlers_deep "echo-starter/internal/services/handlers/deep"
	services_handlers_error "echo-starter/internal/services/handlers/error"
	services_handlers_home "echo-starter/internal/services/handlers/home"

	core_contracts "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/core"
	services_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/claimsprincipal"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Startup struct {
	config *tex_config.Config
	ctrl   *gomock.Controller
}

func NewStartup() *Startup {
	return &Startup{
		config: &tex_config.Config{},
		ctrl:   gomock.NewController(nil),
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

	core_services_timeutils.AddTimeParse(builder)
	services_cookies.AddSingletonISecureCookie(builder)
	// AUTH SERVICES
	//----------------------------------------------------------------------------------------------------------------------
	services_auth_authenticator.AddSingletonIOIDCAuthenticator(builder)

	services_handlers_home.AddScopedIHandler(builder)
	services_handlers_deep.AddScopedIHandler(builder)
	services_handlers_error.AddScopedIHandler(builder)
	services_handlers_about.AddScopedIHandler(builder)

	// AUTH HANDLERS
	//----------------------------------------------------------------------------------------------------------------------
	services_handlers_auth_login.AddScopedIHandler(builder)
	services_handlers_auth_profiles.AddScopedIHandler(builder)
	services_handlers_auth_callback.AddScopedIHandler(builder)
	services_handlers_auth_logout.AddScopedIHandler(builder)
	services_handlers_auth_unauthorized.AddScopedIHandler(builder)

	// ACCOUNT SERVICES
	//----------------------------------------------------------------------------------------------------------------------
	services_handlers_accounts.AddScopedIHandler(builder)
	services_handlers_api_accounts.AddScopedIHandler(builder)

	// ARTISTS CRUD API
	//----------------------------------------------------------------------------------------------------------------------
	services_handlers_artists.AddScopedIHandler(builder)
	services_handlers_api_artists.AddScopedIHandler(builder)
	services_handlers_api_artists_artist.AddScopedIHandler(builder)
	services_handlers_api_artists_artist_albums.AddScopedIHandler(builder)

	services_handler.AddSingletonIHandlerFactory(builder)
	services_core_claimsprincipal.AddScopedIClaimsPrincipal(builder)
	services_claimsprovider.AddSingletonIClaimsProviderMock(builder, s.ctrl)
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
	e.Use(middleware_session.EnsureAuthTokenRefresh(shared.RootContainer))
	e.Use(middleware_claimsprincipal.AuthenticatedSessionToClaimsPrincipalMiddleware(root))
	e.Use(middleware_claimsprincipal.FinalAuthVerificationMiddlewareUsingClaimsMap(echostarter_auth.BuildGrpcEntrypointPermissionsClaimsMap(), true))
	return nil
}
