package contextaccessor

import (
	"github.com/labstack/echo/v4"
)

//go:generate genny -pkg $GOPACKAGE -in=../../../genny/sarulabsdi/interface-types.go -out=gen-$GOFILE gen "InterfaceType=IInternalEchoContextAccessor,IEchoContextAccessor"

//go:generate mockgen -package=$GOPACKAGE -destination=../../mocks/$GOPACKAGE/mock_$GOFILE   echo-starter/internal/contracts/$GOPACKAGE IInternalEchoContextAccessor,IEchoContextAccessor

type (
	// IEchoContextAccessor ...
	IEchoContextAccessor interface {
		GetContext() echo.Context
	}
	IInternalEchoContextAccessor interface {
		IEchoContextAccessor
		SetContext(echo.Context)
	}
)
