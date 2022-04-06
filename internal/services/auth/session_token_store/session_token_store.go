package session_token_store

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	"echo-starter/internal/session"
	"encoding/json"
	"errors"
	"reflect"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_contextaccessor "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/contextaccessor"
	contracts_cookies "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/cookies"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"golang.org/x/oauth2"
)

type (
	service struct {
		Logger              contracts_logger.ILogger                       `inject:""`
		EchoContextAccessor contracts_contextaccessor.IEchoContextAccessor `inject:""`
		SecureCookie        contracts_cookies.ISecureCookie                `inject:""`
		cachedToken         *oauth2.Token
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
	c := s.EchoContextAccessor.GetContext()
	session.TerminateAuthSession(c)
	s.cachedToken = nil
	return nil
}
func (s *service) GetToken() (*oauth2.Token, error) {
	return s.cachedToken, nil
}
func (s *service) SlideOutExpiration() error {
	c := s.EchoContextAccessor.GetContext()
	authSess := session.GetAuthSession(c)
	return authSess.Save(c.Request(), c.Response())
}

func (s *service) GetTokenByIdempotencyKey(bindingKey string) (*oauth2.Token, error) {
	if s.cachedToken == nil {
		c := s.EchoContextAccessor.GetContext()
		authSess := session.GetAuthSession(c)
		mybindingKey, ok := authSess.Values["binding_key"]
		if !ok {
			return nil, errors.New("binding_key not found")
		}
		if mybindingKey.(string) != bindingKey {
			return nil, errors.New("binding_key doesn't match with the one in session")
		}
		authTokens, ok := authSess.Values["tokens"]
		if !ok {
			return nil, errors.New("tokens not found")
		}
		if !ok || core_utils.IsNil(authTokens) {
			return nil, errors.New("tokens not found")
		}
		var token *oauth2.Token = &oauth2.Token{}
		authArtifactsStr := authTokens.(string)
		err := json.Unmarshal([]byte(authArtifactsStr), &token)
		if err != nil {
			return nil, err
		}
		s.cachedToken = token
	}
	return s.cachedToken, nil
}
func (s *service) StoreTokenByIdempotencyKey(bindingKey string, token *oauth2.Token) error {
	c := s.EchoContextAccessor.GetContext()
	authTokensB, err := json.Marshal(token)
	if err != nil {
		return err
	}
	authSess := session.GetAuthSession(c)
	authSess.Values["tokens"] = string(authTokensB)
	authSess.Values["binding_key"] = bindingKey
	err = authSess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}
	s.cachedToken = token
	return nil
}
