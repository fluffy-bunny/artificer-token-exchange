package unauthorized

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	contracts_handler "echo-starter/internal/contracts/handler"
	"echo-starter/internal/templates"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Authenticator   contracts_auth.IOIDCAuthenticator               `inject:"authenticator"`
		ClaimsPrincipal contracts_core_claimsprincipal.IClaimsPrincipal `inject:"claimsPrincipal"`
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
		wellknown.UnauthorizedPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {
	return templates.Render(c, s.ClaimsPrincipal, http.StatusOK, "views/auth/unauthorized/index", map[string]interface{}{
		"claims": s.ClaimsPrincipal.GetClaims(),
	})
}
