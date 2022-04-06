package cookie_token_store

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	"echo-starter/internal/session"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	contracts_config "echo-starter/internal/contracts/config"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_contextaccessor "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/contextaccessor"
	contracts_cookies "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/cookies"
	di "github.com/fluffy-bunny/sarulabsdi"
	"golang.org/x/oauth2"
)

type (
	service struct {
		Config              *contracts_config.Config                       `inject:""`
		Logger              contracts_logger.ILogger                       `inject:""`
		EchoContextAccessor contracts_contextaccessor.IEchoContextAccessor `inject:""`
		SecureCookie        contracts_cookies.ISecureCookie                `inject:""`
		cachedToken         *oauth2.Token
	}
	cookieContainer struct {
		ID    string
		Token *oauth2.Token
	}
)

func assertImplementation() {
	var _ contracts_auth.ITokenStore = (*service)(nil)
	var _ contracts_auth.IInternalTokenStore = (*service)(nil)

}

var reflectType = reflect.TypeOf((*service)(nil))

// AddScopedITokenStore registers the *service as a singleton.
func AddScopedITokenStore(builder *di.Builder) {
	contracts_auth.AddScopedITokenStore(builder, reflectType, contracts_auth.ReflectTypeIInternalTokenStore)
}
func (s *service) Clear() error {
	authCookieName, err := s._getAuthCookieName()
	if err != nil {
		return err
	}
	s.SecureCookie.DeleteCookie(authCookieName)
	return nil
}

func (s *service) _getAuthCookieName() (string, error) {
	c := s.EchoContextAccessor.GetContext()
	sess := session.GetSession(c)

	idempotencyKey, ok := sess.Values["binding_key"]
	if !ok {
		return "", errors.New("idempotency key not found in session")
	}
	authCookieName := fmt.Sprintf("%s_%s", s.Config.AuthCookieName, idempotencyKey)
	return authCookieName, nil
}
func (s *service) SlideOutExpiration() error {
	authCookieName, err := s._getAuthCookieName()
	if err != nil {
		return err
	}
	return s.SecureCookie.RefreshCookie(authCookieName, time.Duration(s.Config.SessionMaxAgeSeconds)*time.Second)
}

func (s *service) GetToken() (*oauth2.Token, error) {
	return s.cachedToken, nil
}
func (s *service) GetTokenByIdempotencyKey(idempotencyKey string) (*oauth2.Token, error) {
	if s.cachedToken == nil {
		authCookieName, err := s._getAuthCookieName()
		if err != nil {
			return nil, err
		}
		jsonS, err := s.SecureCookie.GetCookieValue(authCookieName)
		if err != nil {
			return nil, err
		}
		var container *cookieContainer = &cookieContainer{}
		err = json.Unmarshal([]byte(jsonS), container)
		if err != nil {
			return nil, err
		}
		if container.ID != idempotencyKey {
			s.Logger.Error().Str("request_binding_key", idempotencyKey).
				Str("stored_binding_key", container.ID).Msg("idempotencyKey does not match cookieId")
			return nil, errors.New("binding key requsted doesn't match the one stored")
		}
		s.cachedToken = container.Token
	}
	return s.cachedToken, nil
}
func (s *service) StoreTokenByIdempotencyKey(idempotencyKey string, token *oauth2.Token) error {
	payload := &cookieContainer{
		ID:    idempotencyKey,
		Token: token,
	}
	jsonB, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	expire := time.Now().Add(time.Duration(s.Config.SessionMaxAgeSeconds) * time.Second)
	authCookieName := fmt.Sprintf("%s_%s", s.Config.AuthCookieName, idempotencyKey)
	s.SecureCookie.SetCookieValue(authCookieName, string(jsonB), expire)
	s.cachedToken = token
	return nil
}
