package dev

import (
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"

	contracts_auth "echo-starter/internal/contracts/auth"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"
	contracts_contextaccessor "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/contextaccessor"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger              contracts_logger.ILogger                       `inject:""`
		OIDCAuthenticator   core_contracts_oidc.IOIDCAuthenticator         `inject:""`
		TokenStore          contracts_auth.IInternalTokenStore             `inject:""`
		EchoContextAccessor contracts_contextaccessor.IEchoContextAccessor `inject:""`
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
			contracts_handler.POST,
			contracts_handler.GET,
		},
		wellknown.APIDevPath)
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

type params struct {
	Directive string `param:"directive" query:"directive" header:"directive" form:"directive" json:"directive" xml:"directive"`
}

func (s *service) get(c echo.Context) error {
	u := new(params)
	if err := c.Bind(u); err != nil {
		return err
	}
	switch u.Directive {
	case "session":
		return s.getSessionData(c)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid directive")
	}
}
func (s *service) getSessionData(c echo.Context) error {
	sess := session.GetSession(c)
	sess.Save(c.Request(), c.Response())
	result := make(map[string]string)
	for k, v := range sess.Values {
		key := k.(string)
		switch v.(type) {
		case string:
			result[key] = v.(string)
		case []byte:

			result[key] = string(v.([]byte))
		}
	}
	return c.JSON(http.StatusOK, result)
}
func (s *service) post(c echo.Context) error {
	u := new(params)
	if err := c.Bind(u); err != nil {
		return err
	}
	switch u.Directive {
	case "session":
		return s.getSessionData(c)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid directive")
	}
}
