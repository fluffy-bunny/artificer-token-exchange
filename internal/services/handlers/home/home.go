package home

import (
	"echo-starter/internal/templates"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"

	contracts_config "echo-starter/internal/contracts/config"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_timeutils "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/timeutils"
	contracts_container "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/container"
	contracts_contextaccessor "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/contextaccessor"
	contracts_cookies "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/cookies"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		// Required and Useful services that the runtime registers
		//---------------------------------------------------------------------------------------------
		ContainerAccessor   contracts_container.ContainerAccessor           `inject:""`
		TimeNow             contracts_timeutils.TimeNow                     `inject:""`
		TimeParse           contracts_timeutils.TimeParse                   `inject:""`
		Logger              contracts_logger.ILogger                        `inject:""`
		ClaimsPrincipal     contracts_core_claimsprincipal.IClaimsPrincipal `inject:""`
		SecureCookie        contracts_cookies.ISecureCookie                 `inject:""`
		EchoContextAccessor contracts_contextaccessor.IEchoContextAccessor  `inject:""`
		//---------------------------------------------------------------------------------------------

		// internal services
		Config *contracts_config.Config `inject:""`
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

func (s *service) Ctor() {

}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {
	s.Logger.Info().Str("timeNow", s.TimeNow().String()).Send()
	return templates.Render(c, s.ClaimsPrincipal, http.StatusOK, "views/home/index", map[string]interface{}{})
}
