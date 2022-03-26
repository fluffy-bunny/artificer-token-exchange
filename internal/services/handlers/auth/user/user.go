package logout

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	auth_shared "echo-starter/internal/contracts/auth/shared"
	contracts_handler "echo-starter/internal/contracts/handler"
	"echo-starter/internal/session"
	"echo-starter/internal/utils"
	"echo-starter/internal/wellknown"
	"encoding/json"
	"net/http"
	"reflect"

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
		wellknown.UserPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {
	sess := session.GetSession(c)
	jsonProfile, _ := sess.Values[auth_shared.ProfileSessionKey]
	if jsonProfile != nil {
		var profile map[string]interface{}
		json.Unmarshal(jsonProfile.([]byte), &profile)
		jsonProfileS := utils.PrettyJSON(profile)
		return c.String(http.StatusOK, jsonProfileS)
	} else {
		return c.String(http.StatusOK, "No profile found")
	}
}
