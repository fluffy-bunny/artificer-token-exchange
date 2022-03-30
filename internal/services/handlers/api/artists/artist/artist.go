package artist

import (
	contracts_handler "echo-starter/internal/contracts/handler"
	artists_shared "echo-starter/internal/services/handlers/api/artists/shared"
	"echo-starter/internal/wellknown"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	linq "github.com/ahmetb/go-linq/v3"
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
			contracts_handler.POST,
		},
		wellknown.APIArtistsIdPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type params struct {
	ID string `param:"id" query:"id" header:"id" form:"id" json:"id" xml:"id"`
}

func (s *service) Do(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.get(c)
	case http.MethodPost:
		return s.post(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}
}
func (s *service) post(c echo.Context) error {
	fmt.Println("artist: post")
	decoder := json.NewDecoder(c.Request().Body)
	var body map[string]interface{}
	if err := decoder.Decode(&body); err != nil {
		return err
	}
	fmt.Println(body)
	return c.JSON(http.StatusOK, body)

}
func (s *service) get(c echo.Context) error {

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
	return c.JSON(http.StatusNotFound, artists_shared.Artist{
		Id:   artists[0].Id,
		Name: artists[0].Name,
	})

}
