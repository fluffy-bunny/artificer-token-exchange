package logout

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	contracts_handler "echo-starter/internal/contracts/handler"
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"
	"strings"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	httpVerbs := []contracts_handler.HTTPVERB{
		contracts_handler.GET,
	}
	httpVerbS := []string{}
	for _, httpVerb := range httpVerbs {
		httpVerbS = append(httpVerbS, httpVerb.String())
	}
	path := wellknown.LogoutPath
	metadata := map[string]interface{}{
		"path":      path,
		"httpVerbs": httpVerbs,
	}

	log.Info().
		Str("DI", "IHandler").
		Str("path", path).
		Str("httpVerbs", strings.Join(httpVerbS, "|")).Send()
	contracts_handler.AddScopedIHandlerWithMetadata(builder, reflectType, metadata)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {
	session.TerminateSession(c)

	// Redirect to home page.
	c.Redirect(http.StatusTemporaryRedirect, "/")
	return nil
}
