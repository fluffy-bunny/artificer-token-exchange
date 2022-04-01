package auth

import (
	"golang.org/x/oauth2"
)

//go:generate genny -pkg $GOPACKAGE -in=../../../genny/sarulabsdi/interface-types.go -out=gen-$GOFILE gen "InterfaceType=IAuthArtifacts"

//go:generate mockgen -package=$GOPACKAGE -destination=../../mocks/$GOPACKAGE/mock_$GOFILE   echo-starter/internal/contracts/$GOPACKAGE IAuthArtifacts

type (
	// IAuthArtifacts ...
	IAuthArtifacts interface {
		GetToken() *oauth2.Token
	}
)
