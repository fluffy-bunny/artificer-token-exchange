package auth_artifacts

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	"reflect"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type (
	service struct {
		Logger     contracts_logger.ILogger   `inject:""`
		TokenStore contracts_auth.ITokenStore `inject:""`
	}
)

func assertImplementation() {
	var _ contracts_auth.IAuthArtifacts = (*service)(nil)

}

var reflectType = reflect.TypeOf((*service)(nil))

// AddScopedIAuthArtifacts registers the *service as a singleton.
func AddScopedIAuthArtifacts(builder *di.Builder) {
	log.Info().Str("DI", "IAuthArtifacts - SCOPED").Send()
	contracts_auth.AddScopedIAuthArtifacts(builder, reflectType)
}
func (s *service) GetName() string {
	return "oidc"
}
func (s *service) Refresh() {

}
func (s *service) GetToken() *oauth2.Token {
	token, _ := s.TokenStore.GetToken()
	return token
}
