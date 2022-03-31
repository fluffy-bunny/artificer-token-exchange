package accounts

import (
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type (
	service struct {
		Logger            contracts_logger.ILogger               `inject:""`
		OIDCAuthenticator core_contracts_oidc.IOIDCAuthenticator `inject:""`
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
		},
		wellknown.APIAccountsPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {
	switch c.Request().Method {

	case http.MethodPost:
		return s.post(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}
}

type postParams struct {
	Directive string `param:"directive" query:"directive" header:"directive" form:"directive" json:"directive" xml:"directive"`
}

func (s *service) post(c echo.Context) error {
	u := new(postParams)
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
func (s *service) postForceRefresh(c echo.Context) error {
	ctx := c.Request().Context()
	fmt.Println("accounts: postForceRefresh")
	authSession := session.GetAuthSession(c)
	var denied = func() error {
		return echo.NewHTTPError(http.StatusForbidden, "not authorized")
	}
	if authSession == nil {
		return denied()
	}

	for {
		authArtifacts, ok := authSession.Values["tokens"]
		if !ok || core_utils.IsNil(authArtifacts) {
			break
		}
		var token *oauth2.Token = &oauth2.Token{}
		authArtifactsStr := authArtifacts.(string)
		err := json.Unmarshal([]byte(authArtifactsStr), &token)
		if err != nil {
			log.Error().Err(err).Msg("unmarshal token")
			break
		}
		// make the token expired so that tokenSource will refresh it
		token.Expiry = time.Now().Add(time.Duration(-60) * time.Second)
		tokenSource := s.OIDCAuthenticator.GetTokenSource(ctx, token)
		// token source will not do the refresh for us
		newToken, err := tokenSource.Token()
		if err != nil {
			log.Warn().Err(err).Msg("refresh token failed")
			break
		}
		if newToken.AccessToken != token.AccessToken {
			// we need to save this one into the session
			authTokensB, err := json.Marshal(token)
			if err == nil {
				authSession.Values["tokens"] = string(authTokensB)
			}
			authSession.Save(c.Request(), c.Response())
		}
		return c.JSON(http.StatusOK, "ok")
	}
	return denied()
}
