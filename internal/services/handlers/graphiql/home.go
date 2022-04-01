package graphiql

import (
	"echo-starter/internal/templates"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"

	contracts_config "echo-starter/internal/contracts/config"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_container "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/container"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Config            *contracts_config.Config                        `inject:""`
		Logger            contracts_logger.ILogger                        `inject:""`
		ContainerAccessor contracts_container.ContainerAccessor           `inject:""`
		ClaimsPrincipal   contracts_core_claimsprincipal.IClaimsPrincipal `inject:""`
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
		wellknown.GraphiQLPath)
}

func (s *service) Ctor() {

}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {

	return templates.Render(c, s.ClaimsPrincipal, http.StatusOK, "views/graphiql/index", map[string]interface{}{})
}
