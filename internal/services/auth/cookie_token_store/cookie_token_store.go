package cookie_token_store

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	"reflect"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_contextaccessor "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/contextaccessor"
	contracts_cookies "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/cookies"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
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
}
func (s *service) Clear() error {

	return nil
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddScopedITokenStore registers the *service as a singleton.
func AddScopedITokenStore(builder *di.Builder) {
	log.Info().Str("DI", "ITokenStore,IInternalTokenStore - COOKIE SCOPED").Send()
	contracts_auth.AddScopedITokenStore(builder, reflectType, contracts_auth.ReflectTypeIInternalTokenStore)
}
func (s *service) GetToken() (*oauth2.Token, error) {
	return s.cachedToken, nil
}
func (s *service) GetTokenByIdompotencyKey(idompotencyKey string) (*oauth2.Token, error) {
	return nil, nil
}
func (s *service) StoreTokenByIdompotencyKey(idompotencyKey string, token *oauth2.Token) error {
	return nil
}
