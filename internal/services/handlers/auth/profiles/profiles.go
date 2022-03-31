package profiles

import (
	contracts_claimsprovider "echo-starter/internal/contracts/claimsprovider"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"

	"echo-starter/internal/templates"
	"errors"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger          contracts_logger.ILogger                        `inject:""`
		ClaimsProvider  contracts_claimsprovider.IClaimsProvider        `inject:""`
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
			contracts_handler.POST,
		},
		wellknown.ProfilesPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
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

type auth struct {
	CSRF string `param:"csrf" query:"csrf" header:"csrf" form:"csrf" json:"csrf" xml:"csrf"`
}

func (s *service) getUserId(c echo.Context) (string, error) {
	subClaims := s.ClaimsPrincipal.GetClaimsByType("sub")
	if len(subClaims) == 0 {
		return "", errors.New("no sub claim found")
	}
	userId := subClaims[0].Value
	return userId, nil
}

func (s *service) get(c echo.Context) error {
	userId, err := s.getUserId(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/error?message=no+sub+claim+found")
	}
	profiles, err := s.ClaimsProvider.GetProfiles(userId)
	if err != nil {
		c.Redirect(http.StatusFound, "/error?message=error+retrieving+profiles")

	}
	return templates.Render(c, s.ClaimsPrincipal, http.StatusOK, "views/auth/profiles/index", map[string]interface{}{
		"profiles": profiles,
	})

}

type postParams struct {
	Profile string `param:"profile" query:"profile" header:"profile" form:"profile" json:"profile" xml:"profile"`
}

func (s *service) post(c echo.Context) error {
	u := new(postParams)
	if err := c.Bind(u); err != nil {
		return err
	}
	c.Logger().Debug(u)

	return c.Redirect(http.StatusFound, "/")
}
