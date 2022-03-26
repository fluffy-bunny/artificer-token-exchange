package deep

import (
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
		wellknown.DeepPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type deepParams struct {
	ID   string `param:"id" query:"id" header:"id" form:"id" json:"id" xml:"id"`
	Name string `param:"name" query:"name" header:"name" form:"name" json:"name" xml:"name"`
}

func (s *service) Do(c echo.Context) error {

	u := new(deepParams)
	if err := c.Bind(u); err != nil {
		return err
	}

	return templates.Render(c, s.ClaimsPrincipal, http.StatusOK, "views/deep/deep", map[string]interface{}{
		"params": u,
	})

}
