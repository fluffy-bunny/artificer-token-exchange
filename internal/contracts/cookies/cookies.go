package cookies

import (
	"time"

	"github.com/labstack/echo/v4"
)

//go:generate genny -pkg $GOPACKAGE -in=../../../genny/sarulabsdi/interface-types.go -out=gen-$GOFILE gen "InterfaceType=ISecureCookie"

//go:generate mockgen -package=$GOPACKAGE -destination=../../mocks/$GOPACKAGE/mock_$GOFILE   echo-starter/internal/contracts/$GOPACKAGE ISecureCookie

type (
	ISecureCookie interface {
		SetCookieValue(c echo.Context, name string, value string, expires time.Time) error
		GetCookieValue(c echo.Context, name string) (string, error)
		DeleteCookie(c echo.Context, name string) error
		RefreshCookie(c echo.Context, name string, durration time.Duration) error
	}
)
