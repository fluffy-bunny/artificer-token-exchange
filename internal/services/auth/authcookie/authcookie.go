package authcookie

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	contracts_config "echo-starter/internal/contracts/config"
	"net/http"
	"reflect"
	"time"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type (
	service struct {
		AppConfig *contracts_config.Config `inject:"config"`
	}
)

func assertImplementation() {
	var _ contracts_auth.IAuthCookie = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddSingletonIAuthCookie registers the *service as a singleton.
func AddSingletonIAuthCookie(builder *di.Builder) {
	log.Info().Str("DI", "IAuthCookie").Send()
	contracts_auth.AddSingletonIAuthCookie(builder, reflectType)
}

func (s *service) SetAuthCookieValue(c echo.Context, value string) error {
	cookie := new(http.Cookie)
	cookie.Name = s.AppConfig.AuthCookieName
	cookie.Value = value
	cookie.Expires = time.Now().Add(time.Duration(s.AppConfig.AuthCookieExpireSeconds) * time.Second)
	c.SetCookie(cookie)
	return nil
}
func (s *service) GetAuthCookieValue(c echo.Context) (string, error) {
	cookie, err := c.Cookie(s.AppConfig.AuthCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
func (s *service) DeleteAuthCookie(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = s.AppConfig.AuthCookieName
	cookie.Value = ""
	cookie.Expires = time.Now()
	c.SetCookie(cookie)
	return nil
}
func (s *service) RefreshAuthCookie(c echo.Context) error {
	cookie, err := c.Cookie(s.AppConfig.AuthCookieName)
	if err != nil {
		return err
	}
	cookie.Expires = time.Now().Add(time.Duration(s.AppConfig.AuthCookieExpireSeconds) * time.Second)
	c.SetCookie(cookie)
	return nil
}
