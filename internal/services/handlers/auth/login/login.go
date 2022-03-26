package login

import (
	"crypto/rand"
	contracts_auth "echo-starter/internal/contracts/auth"
	auth_shared "echo-starter/internal/contracts/auth/shared"
	contracts_handler "echo-starter/internal/contracts/handler"
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"

	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Authenticator contracts_auth.IOIDCAuthenticator `inject:"authenticator"`
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
	url := s.Authenticator.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
	return nil
}
func (s *service) generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
