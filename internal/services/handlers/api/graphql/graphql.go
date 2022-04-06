package graphql

import (
	"bytes"
	contracts_auth "echo-starter/internal/contracts/auth"
	contracts_config "echo-starter/internal/contracts/config"
	"echo-starter/internal/wellknown"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger     contracts_logger.ILogger   `inject:""`
		TokenStore contracts_auth.ITokenStore `inject:""`
		Config     *contracts_config.Config   `inject:""`
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
		wellknown.GraphQLEndpointPath)
}

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

func (s *service) post(c echo.Context) error {
	token, err := s.TokenStore.GetToken()
	if err != nil || token == nil {
		return echo.NewHTTPError(http.StatusForbidden, "not authorized")
	}
	if core_utils.IsEmptyOrNil(token.AccessToken) {
		s.Logger.Error().Msg("access token is empty")
		return echo.NewHTTPError(http.StatusForbidden, "not authorized")
	}

	b, err := ioutil.ReadAll(c.Request().Body)

	url := s.Config.GraphQLEndpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		s.Logger.Error().Err(err).Msg("http.NewRequest")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		s.Logger.Error().Err(err).Msg("client.Do")
		return err
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error().Err(err).Msg("ioutil.ReadAll")
		return err
	}
	return c.String(http.StatusOK, string(b))
}
