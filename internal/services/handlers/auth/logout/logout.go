package logout

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	contracts_handler "echo-starter/internal/contracts/handler"
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Authenticator contracts_auth.IOIDCAuthenticator `inject:"authenticator"`
		AuthCookie    contracts_auth.IAuthCookie        `inject:""`
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
		wellknown.LogoutPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {
	// TODO in larger systems there can be a session that holds may users, think a netflix profile, etc.
	// the profile (or baby user) is removed fro the session vs the entire session
	session.TerminateSession(c)

	s.AuthCookie.DeleteAuthCookie(c)

	// Redirect to home page.
	c.Redirect(http.StatusFound, "/")
	return nil
}
