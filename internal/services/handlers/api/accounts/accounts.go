package accounts

import (
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"fmt"
	"net/http"
	"reflect"
	"time"

	contracts_auth "echo-starter/internal/contracts/auth"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_contracts_oauth2 "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oauth2"

	contracts_contextaccessor "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/contextaccessor"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type (
	service struct {
		Logger              contracts_logger.ILogger                       `inject:""`
		OAuth2Authenticator core_contracts_oauth2.IOAuth2Authenticator     `inject:""`
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
		wellknown.APIAccountsPath)
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

	return c.JSON(http.StatusOK, sess.Values)
}
func (s *service) post(c echo.Context) error {
	u := new(params)
	if err := c.Bind(u); err != nil {
		return err
	}
	switch u.Directive {
	case "force-refresh":
		return s.postForceRefresh(c)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid directive")
	}
}
func (s *service) _getbindingKey(c echo.Context) (string, error) {
	sess := session.GetSession(c)

	bindingKey, ok := sess.Values["binding_key"]
	if !ok {
		return "", fmt.Errorf("binding_key not found")
	}
	return bindingKey.(string), nil
}
func (s *service) postForceRefresh(c echo.Context) error {
	ctx := c.Request().Context()
	fmt.Println("accounts: postForceRefresh")

	var denied = func() error {
		return echo.NewHTTPError(http.StatusForbidden, "not authorized")
	}

	bindingKey, err := s._getbindingKey(c)
	if err != nil {
		s.Logger.Error().Err(err).Msg("binding_key not found")
		return denied()
	}

	for {
		token, err := s.TokenStore.GetTokenByIdempotencyKey(bindingKey)
		if err != nil {
			s.Logger.Error().Msg("TokenStore.GetTokenBybindingKey failed")
			break
		}

		// make the token expired so that tokenSource will refresh it
		token.Expiry = time.Now().Add(time.Duration(-60) * time.Second)
		tokenSource := s.OAuth2Authenticator.GetTokenSource(ctx, token)

		// token source will not do the refresh for us
		newToken, err := tokenSource.Token()
		if err != nil {
			log.Warn().Err(err).Msg("refresh token failed")
			return c.JSON(http.StatusOK, err.Error())
		}
		if newToken.AccessToken != token.AccessToken {
			err = s.TokenStore.StoreTokenByIdempotencyKey(bindingKey, newToken)
			if err != nil {
				s.Logger.Error().Err(err).Msg("TokenStore.StoreTokenBybindingKey failed")
				break
			}
		}

		return c.JSON(http.StatusOK, "ok")
	}
	return denied()
}
