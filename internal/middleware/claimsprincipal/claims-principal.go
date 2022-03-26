package claimsprincipal

import (
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	echostarter_wellknown "echo-starter/internal/wellknown"
	"net/http"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	middleware_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/claimsprincipal"
	middleware_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/oidc"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func recursiveAddClaim(claimsConfig *middleware_oidc.ClaimsConfig, claimsPrincipal contracts_core_claimsprincipal.IClaimsPrincipal) {
	for _, claimFact := range claimsConfig.AND {
		claimsPrincipal.AddClaim(claimFact.Claim)
	}
	for _, claimFact := range claimsConfig.OR {
		claimsPrincipal.AddClaim(claimFact.Claim)
	}
	if claimsConfig.Child != nil {
		recursiveAddClaim(claimsConfig.Child, claimsPrincipal)
	}
}

// DevelopmentMiddlewareUsingClaimsMap use this in development if you are making an api only service
// it literally just adds the claims to the principal that the api demands it has to be authorized.
func DevelopmentMiddlewareUsingClaimsMap(entrypointClaimsMap map[string]*middleware_oidc.EntryPointConfig, enableZeroTrust bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			scopedContainer := c.Get(echostarter_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)
			elem, ok := entrypointClaimsMap[c.Path()]
			if ok {
				recursiveAddClaim(elem.ClaimsConfig, claimsPrincipal)
			}
			return next(c)
		}
	}
}

type OnUnauthorizedAction int64

const (
	OnUnauthorizedAction_Unspecified OnUnauthorizedAction = 0
	OnUnauthorizedAction_Redirect                         = 1
)

type EntryPointConfigEx struct {
	middleware_oidc.EntryPointConfig
	OnUnauthorizedAction OnUnauthorizedAction
}

func AuthenticatedSessionToClaimsPrincipalMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.GetSession(c)
			_, ok := sess.Values[wellknown.ClaimTypeAuthenticated]
			if ok {
				tryAddProfileClaim := func(sessVItem string, claimsPrincipal contracts_core_claimsprincipal.IClaimsPrincipal) {
					pItem, ok := sess.Values[sessVItem]
					if ok {
						claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
							Type:  sessVItem,
							Value: pItem.(string)})
					}
				}
				scopedContainer := c.Get(echostarter_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
				claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)

				claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
					Type:  echostarter_wellknown.ClaimTypeAuthenticated,
					Value: "*"})

				tryAddProfileClaim("id:sub", claimsPrincipal)
				tryAddProfileClaim("id:name", claimsPrincipal)
				tryAddProfileClaim("id:email", claimsPrincipal)
				tryAddProfileClaim("id:iss", claimsPrincipal)
			}
			return next(c)
		}
	}
}
func FinalAuthVerificationMiddlewareUsingClaimsMap(entrypointClaimsMap map[string]*middleware_oidc.EntryPointConfig, enableZeroTrust bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			subLogger := log.With().
				Bool("enableZeroTrust", enableZeroTrust).
				Str("FullMethod", c.Path()).
				Logger()
			debugEvent := subLogger.Debug()

			scopedContainer := c.Get(echostarter_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)

			authenticated := claimsPrincipal.HasClaimType(wellknown.ClaimTypeAuthenticated)
			elem, ok := entrypointClaimsMap[c.Path()]
			permissionDeniedFunc := func() error {
				if !authenticated {
					directive, ok := elem.MetaData["onUnauthenticated"]
					if ok && directive == "login" {
						return c.Redirect(http.StatusFound, "/login?redirect_url="+c.Request().URL.String())
					}
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			if !ok && enableZeroTrust {
				debugEvent.Msg("FullMethod not found in entrypoint claims map")
				return permissionDeniedFunc()
			}
			if !middleware_claimsprincipal.Validate(debugEvent, elem.ClaimsConfig, claimsPrincipal) {
				debugEvent.Msg("ClaimsConfig validation failed")
				return permissionDeniedFunc()
			}
			return next(c)
		}
	}
}
