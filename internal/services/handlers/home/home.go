package home

import (
	contracts_handler "echo-starter/internal/contracts/handler"

	"echo-starter/internal/wellknown"

	"reflect"

	"echo-starter/internal/models"
	"net/http"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	di "github.com/fluffy-bunny/sarulabsdi"

	"github.com/labstack/echo/v4"
)

type (
	service struct {
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
		wellknown.HomePath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {

	return c.Render(http.StatusOK, "content", map[string]interface{}{
		"isAuthenticated": func() bool { return s.ClaimsPrincipal.HasClaimType(wellknown.ClaimTypeAuthenticated) },
		"paths":           models.NewPaths(),
		"claims":          s.ClaimsPrincipal.GetClaims(),
	})
}
