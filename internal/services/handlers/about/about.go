package about

import (
	"echo-starter/internal/templates"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"
	"strings"

	golinq "github.com/ahmetb/go-linq/v3"
	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_container "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/container"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger            contracts_logger.ILogger                        `inject:""`
		ClaimsPrincipal   contracts_core_claimsprincipal.IClaimsPrincipal `inject:""`
		ContainerAccessor contracts_container.ContainerAccessor           `inject:""`
		HandlerFactory    contracts_handler.IHandlerFactory               `inject:""`
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
		wellknown.AboutPath)
}

func (s *service) Ctor() {

}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {
	handlerDefinitions := contracts_handler.GetIHandlerDefinitions(s.ContainerAccessor())
	s.Logger.Info().Msg("about")
	type row struct {
		Verbs string
		Path  string
	}

	var rows []row

	golinq.From(handlerDefinitions).Select(func(c interface{}) interface{} {
		def := c.(*di.Def)
		path := def.MetaData["path"].(string)
		httpVerbs, _ := def.MetaData["httpVerbs"].([]contracts_handler.HTTPVERB)
		verbBldr := strings.Builder{}

		for idx, verb := range httpVerbs {
			verbBldr.WriteString(verb.String())
			if idx < len(httpVerbs)-1 {
				verbBldr.WriteString(",")
			}
		}
		return row{
			Verbs: verbBldr.String(),
			Path:  path,
		}

	}).OrderBy(func(i interface{}) interface{} {
		return i.(row).Path
	}).ToSlice(&rows)

	return templates.Render(c, s.ClaimsPrincipal, http.StatusOK, "views/about/index", map[string]interface{}{
		"defs": rows,
	})
}
