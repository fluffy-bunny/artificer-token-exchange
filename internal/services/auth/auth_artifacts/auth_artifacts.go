package auth_artifacts

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	"echo-starter/internal/session"
	"encoding/json"
	"reflect"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_contextaccessor "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/contextaccessor"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type (
	service struct {
		Logger              contracts_logger.ILogger                       `inject:""`
		EchoContextAccessor contracts_contextaccessor.IEchoContextAccessor `inject:""`
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
func (s *service) GetToken() *oauth2.Token {
	sess := session.GetAuthSession(s.EchoContextAccessor.GetContext())
	if sess == nil {
		return nil
	}
	authArtifacts, ok := sess.Values["tokens"]
	if !ok || core_utils.IsNil(authArtifacts) {
		return nil
	}
	var token *oauth2.Token = &oauth2.Token{}
	authArtifactsStr := authArtifacts.(string)
	err := json.Unmarshal([]byte(authArtifactsStr), &token)
	if err != nil {
		s.Logger.Error().Err(err).Msg("unmarshal token")
		return nil

	}
	return token
}
