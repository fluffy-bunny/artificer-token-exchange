package authenticator

import (
	"context"
	contracts_auth "echo-starter/internal/contracts/auth"
	contracts_config "echo-starter/internal/contracts/config"
	"errors"
	"reflect"

	services_oidc "echo-starter/internal/services/oidc"

	"github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type (
	service struct {
		*oidc.Provider
		oauth2.Config
		AppConfig      *contracts_config.Config `inject:"config"`
		oidcProviderEx *services_oidc.Provider
		issuer         string
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
	s.issuer = "https://" + s.AppConfig.Oidc.Domain + "/"
	oidcProviderEx, err := services_oidc.NewProvider(context.Background(), s.issuer)
	if err != nil {
		panic(err)
	}
	s.oidcProviderEx = oidcProviderEx
	provider, err := oidc.NewProvider(
		context.Background(),
		s.issuer,
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

func (s *service) ValidateJWTAccessToken(accessToken string) (*services_oidc.AccessToken, error) {
	verifier := services_oidc.NewJWTAccessTokenVerifier(s.issuer, s.oidcProviderEx.GetRemoteKeySet(), &oidc.Config{
		SkipClientIDCheck: true,
	})
	return verifier.Verify(context.Background(), accessToken)
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

func (s *service) GetTokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	ts := s.Config.TokenSource(ctx, token)
	return ts
}
