package authenticator

import (
	"context"
	contracts_auth "echo-starter/internal/contracts/auth"
	contracts_config "echo-starter/internal/contracts/config"
	"errors"
	"reflect"

	"github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type (
	service struct {
		*oidc.Provider
		oauth2.Config
		AppConfig *contracts_config.Config `inject:"config"`
	}
)

func assertImplementation() {
	var _ contracts_auth.IOIDCAuthenticator = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddSingletonIOIDCAuthenticator registers the *service as a singleton.
func AddSingletonIOIDCAuthenticator(builder *di.Builder) {
	log.Info().Str("DI", "IOIDCAuthenticator").Send()
	contracts_auth.AddSingletonIOIDCAuthenticator(builder, reflectType)
}
func (s *service) Ctor() {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+s.AppConfig.Oidc.Domain+"/",
	)
	if err != nil {
		panic(err)
	}
	s.Provider = provider

	conf := oauth2.Config{
		ClientID:     s.AppConfig.Oidc.ClientID,
		ClientSecret: s.AppConfig.Oidc.ClientSecret,
		RedirectURL:  s.AppConfig.Oidc.CallbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile"},
	}
	s.Config = conf
}
func (s *service) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: s.ClientID,
	}

	return s.Verifier(oidcConfig).Verify(ctx, rawIDToken)

}
