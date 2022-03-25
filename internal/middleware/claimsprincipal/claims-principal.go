package claimsprincipal

import (
	echostarter_wellknown "echo-starter/internal/wellknown"
	"net/http"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	middleware_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/oidc"

	middleware_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/claimsprincipal"
	"github.com/rs/zerolog/log"

	di "github.com/fluffy-bunny/sarulabsdi"

	"github.com/labstack/echo/v4"
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

			permissionDeniedFunc := func() error {
				c.Response().Status = http.StatusUnauthorized
				c.Response().Write([]byte("Permission Denied"))
				return nil
			}
			elem, ok := entrypointClaimsMap[c.Path()]
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
