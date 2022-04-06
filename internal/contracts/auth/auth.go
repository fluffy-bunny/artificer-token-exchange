package auth

import (
	"golang.org/x/oauth2"
)

//go:generate genny -pkg $GOPACKAGE -in=../../../genny/sarulabsdi/interface-types.go -out=gen-$GOFILE gen "InterfaceType=IInternalTokenStore,ITokenStore"

//go:generate mockgen -package=$GOPACKAGE -destination=../../mocks/$GOPACKAGE/mock_$GOFILE   echo-starter/internal/contracts/$GOPACKAGE IInternalTokenStore,ITokenStore

type (

	// ITokenStore is a SCOPED store so nothing global
	ITokenStore interface {
		GetToken() (*oauth2.Token, error)
		Clear() error
	}
	IInternalTokenStore interface {
		ITokenStore
		GetTokenByIdempotencyKey(bindingKey string) (*oauth2.Token, error)
		StoreTokenByIdempotencyKey(bindingKey string, token *oauth2.Token) error
		SlideOutExpiration() error
	}
)
