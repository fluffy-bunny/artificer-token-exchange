package claimsprincipal

import (
	"echo-starter/internal/session"

	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"
	core_echo "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"

	contracts_auth "echo-starter/internal/contracts/auth"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	middleware_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/oidc"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"gopkg.in/square/go-jose.v2/jwt"
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

type OnUnauthorizedAction int64

const (
	OnUnauthorizedAction_Unspecified OnUnauthorizedAction = 0
	OnUnauthorizedAction_Redirect                         = 1
)

type EntryPointConfigEx struct {
	middleware_oidc.EntryPointConfig
	OnUnauthorizedAction OnUnauthorizedAction
}

const middlewareLogName = "token-to-claims-principal"

func AuthenticatedSessionToClaimsPrincipalMiddleware(root di.Container) echo.MiddlewareFunc {
	// get authCookie service once during configuration

	oidcAuthenticator, _ := core_contracts_oidc.SafeGetIOIDCAuthenticatorFromContainer(root)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			for {

				// Skip this if we see an authorization header
				// important: The CSRF middleware is skipped as well if there is an Authorization header
				// So if we get here then we can't be adding any claims if someone got our session
				// always use the HasWellknownAuthHeaders centralized func
				if core_echo.HasWellknownAuthHeaders(c) {
					// this is a cookie/session claims maker so if another authorization scheme is used we will not contribute
					break
				}
				// 1. get our idompontent session
				sess := session.GetSession(c)
				bindingKey, ok := sess.Values["binding_key"]
				if !ok {
					// if we don't  have this the user hasn't logged in
					break
				}
				var terminateAuthSession = func() {
					session.TerminateSession(c)
				}

				scopedContainer := c.Get(core_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
				logger := contracts_logger.GetILoggerFromContainer(scopedContainer)
				errorEvent := logger.GetLogger().Error().Str("middleware", middlewareLogName)
				debugEvent := logger.GetLogger().Debug().Str("middleware", middlewareLogName)

				tokenStore := contracts_auth.GetIInternalTokenStoreFromContainer(scopedContainer)

				token, err := tokenStore.GetTokenByIdempotencyKey(bindingKey.(string))
				if err != nil {
					// not necessarily an error. The tokens could have been removed and our idompotent key could be stale
					debugEvent.Err(err).Msg("Failed to get token")
					terminateAuthSession()
					break
				}

				claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)

				if ok, _ := isJWT(token.AccessToken); ok {
					if oidcAuthenticator != nil {
						accessToken, err := oidcAuthenticator.ValidateJWTAccessToken(token.AccessToken)
						if err != nil {
							errorEvent.Err(err).Msg("ValidateJWTAccessToken failed")
							terminateAuthSession()
							break
						}
						accessTokenClaims := accessToken.ToClaims()
						for _, claim := range accessTokenClaims {
							claimsPrincipal.AddClaim(*claim)
						}

					}
				}

				claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
					Type:  core_wellknown.ClaimTypeAuthenticated,
					Value: "*"})

				break
			}

			return next(c)
		}
	}
}
func isJWT(token string) (bool, error) {
	_, err := jwt.ParseSigned(token)
	if err != nil {
		return false, err
	}
	return true, nil
}
