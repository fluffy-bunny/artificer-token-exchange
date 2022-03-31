package claimsprincipal

import (
	"echo-starter/internal/session"
	"encoding/json"

	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"

	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
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

	authenticator := core_contracts_oidc.GetIOIDCAuthenticatorFromContainer(root)
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
				scopedContainer := c.Get(core_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
				claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)
				for _, claim := range accessTokenClaims {
					claimsPrincipal.AddClaim(*claim)
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
