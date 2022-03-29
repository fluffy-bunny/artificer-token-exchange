package claimsprincipal

import (
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	echostarter_wellknown "echo-starter/internal/wellknown"
	"encoding/json"
	"net/http"

	contracts_auth "echo-starter/internal/contracts/auth"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	middleware_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/claimsprincipal"
	middleware_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/oidc"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
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

func AuthenticatedSessionToClaimsPrincipalMiddleware(root di.Container) echo.MiddlewareFunc {
	// get authCookie service once during configuration

	authenticator := contracts_auth.GetIOIDCAuthenticatorFromContainer(root)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for {
				sess := session.GetAuthSession(c)
				var terminateAuthSession = func() {
					sess.Values = make(map[interface{}]interface{})
					sess.Save(c.Request(), c.Response())
				}

				authArtifacts, ok := sess.Values["tokens"]
				if !ok || core_utils.IsNil(authArtifacts) {
					break
				}
				var token *oauth2.Token = &oauth2.Token{}
				authArtifactsStr := authArtifacts.(string)
				err := json.Unmarshal([]byte(authArtifactsStr), &token)
				if err != nil {
					log.Error().Err(err).Msg("unmarshal token")
					terminateAuthSession()
					break
				}
				accessToken, err := authenticator.ValidateJWTAccessToken(token.AccessToken)
				if err != nil {
					log.Error().Err(err).Msg("ValidateJWTAccessToken failed")
					terminateAuthSession()
					break
				}

				accessTokenClaims := accessToken.ToClaims()
				scopedContainer := c.Get(echostarter_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
				claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)
				for _, claim := range accessTokenClaims {
					claimsPrincipal.AddClaim(*claim)
				}

				claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
					Type:  echostarter_wellknown.ClaimTypeAuthenticated,
					Value: "*"})

				break
			}

			return next(c)
		}
	}
}
func FinalAuthVerificationMiddlewareUsingClaimsMap(entrypointClaimsMap map[string]*middleware_oidc.EntryPointConfig, enableZeroTrust bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			path := c.Path()
			subLogger := log.With().
				Bool("enableZeroTrust", enableZeroTrust).
				Str("FullMethod", path).
				Logger()
			debugEvent := subLogger.Debug()

			scopedContainer := c.Get(echostarter_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)

			authenticated := claimsPrincipal.HasClaimType(wellknown.ClaimTypeAuthenticated)
			elem, ok := entrypointClaimsMap[path]
			permissionDeniedFunc := func() error {
				if !authenticated {
					if !core_utils.IsNil(elem) {
						directive, ok := elem.MetaData["onUnauthenticated"]
						if ok && directive == "login" {
							return c.Redirect(http.StatusFound, "/login?redirect_url="+c.Request().URL.String())
						}
						return c.String(http.StatusUnauthorized, "Unauthorized")
					}
				}
				return c.Redirect(http.StatusFound, "/unauthorized")
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
