package albumns

import (
	artists_shared "echo-starter/internal/services/handlers/api/artists/shared"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"

	linq "github.com/ahmetb/go-linq/v3"
	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger          contracts_logger.ILogger                        `inject:""`
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
		wellknown.APIArtistsIdAlbumsPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type params struct {
	ID string `param:"id" query:"id" header:"id" form:"id" json:"id" xml:"id"`
}

func (s *service) Do(c echo.Context) error {
	u := new(params)
	if err := c.Bind(u); err != nil {
		return err
	}

	var artists []artists_shared.Artist

	linq.From(artists_shared.Artists).Where(func(c interface{}) bool {
		return c.(artists_shared.Artist).Id == u.ID
	}).Select(func(c interface{}) interface{} {
		return c
	}).ToSlice(&artists)

	if len(artists) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusNotFound, artists[0].Albums)

}
