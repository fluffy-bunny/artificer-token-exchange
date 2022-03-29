package auth

import (
	"context"

	services_oidc "echo-starter/internal/services/oidc"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

//go:generate genny -pkg $GOPACKAGE -in=../../../genny/sarulabsdi/interface-types.go -out=gen-$GOFILE gen "InterfaceType=IOIDCAuthenticator"

//go:generate mockgen -package=$GOPACKAGE -destination=../../mocks/$GOPACKAGE/mock_$GOFILE   echo-starter/internal/contracts/$GOPACKAGE IOIDCAuthenticator

type (
	// IOIDCAuthenticator ...
	IOIDCAuthenticator interface {
		GetTokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource
		VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
		ValidateJWTAccessToken(accessToken string) (*services_oidc.AccessToken, error)
		AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
		Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	}
)
