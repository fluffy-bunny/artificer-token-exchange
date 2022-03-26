package auth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

//go:generate genny -pkg $GOPACKAGE -in=../../../genny/sarulabsdi/interface-types.go -out=gen-$GOFILE gen "InterfaceType=IOIDCAuthenticator"

//go:generate mockgen -package=$GOPACKAGE -destination=../../mocks/$GOPACKAGE/mock_$GOFILE   echo-starter/internal/contracts/$GOPACKAGE IOIDCAuthenticator

type (
	// IOIDCAuthenticator ...
	IOIDCAuthenticator interface {
		VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
		AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
		Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	}
)
