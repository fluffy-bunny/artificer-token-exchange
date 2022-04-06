package logout

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	contracts_config "echo-starter/internal/contracts/config"
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"fmt"
	"net/http"
	"reflect"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger          contracts_logger.ILogger                        `inject:""`
		Config          *contracts_config.Config                        `inject:""`
		TokenStore      contracts_auth.ITokenStore                      `inject:""`
		ClaimsPrincipal contracts_core_claimsprincipal.IClaimsPrincipal `inject:""`
	}
)

func assertImplementation() {
	var _ contracts_handler.IHandler = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder *di.Builder) {
	contracts_handler.AddScopedIHandlerEx(builder,
		reflectType,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown.LogoutPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func getMyRootPath(c echo.Context) string {
	return fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host)
}
func (s *service) Do(c echo.Context) error {
	// 1. Clear our auth tokens first.  The middelware can recover if the main session is not cleared
	s.TokenStore.Clear()
	session.TerminateSession(c)

	// 2. log the user out of everything use the auth0 logout endpoint

	// https://auth0.com/docs/api/authentication?javascript#logout
	//GET https://YOUR_DOMAIN/v2/logout?client_id=YOUR_CLIENT_ID&returnTo=LOGOUT_URL
	logoutUrl := fmt.Sprintf("https://%s/v2/logout?client_id=%s&returnTo=%s/",
		s.Config.OIDC.Domain,
		s.Config.OIDC.ClientID,
		getMyRootPath(c))

	// Redirect to home page.
	return c.Redirect(http.StatusFound, logoutUrl)
}
