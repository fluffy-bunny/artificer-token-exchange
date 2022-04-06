package login

import (
	auth_shared "echo-starter/internal/contracts/auth/shared"
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"encoding/json"
	"net/http"
	"reflect"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"
	"github.com/rs/xid"
	"golang.org/x/oauth2"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger          contracts_logger.ILogger                        `inject:""`
		Authenticator   core_contracts_oidc.IOIDCAuthenticator          `inject:""`
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
		wellknown.LoginPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {
	u := new(auth_shared.LoginParms)
	if err := c.Bind(u); err != nil {
		return err
	}
	if core_utils.IsEmptyOrNil(u.RedirectURL) {
		u.RedirectURL = "/"
	}
	jsonLoginParams, err := json.Marshal(u)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	state, err := s.generateRandomState()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())

	}
	sess := session.GetSession(c)
	sess.Values[auth_shared.AuthStateSessionKey] = state
	sess.Values[auth_shared.LoginParamsSessionKey] = jsonLoginParams

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	url := s.Authenticator.AuthCodeURL(state, oauth2.AccessTypeOffline)
	/*
		return templates.Render(c, s.ClaimsPrincipal, http.StatusOK, "views/auth/login/index", map[string]interface{}{
			"url": url,
		})
	*/
	return c.Redirect(http.StatusFound, url)
}
func (s *service) generateRandomState() (string, error) {
	return xid.New().String(), nil
}
