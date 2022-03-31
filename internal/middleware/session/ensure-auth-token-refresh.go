package session

import (
	"context"
	"echo-starter/internal/session"
	"encoding/json"

	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func EnsureAuthTokenRefresh(container di.Container) echo.MiddlewareFunc {
	authenticator := core_contracts_oidc.GetIOIDCAuthenticatorFromContainer(container)
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			authSession := session.GetAuthSession(c)
			if !authSession.IsNew {
				for {
					var terminateAuthSession = func() {
						authSession.Values = make(map[interface{}]interface{})
						authSession.Save(c.Request(), c.Response())
					}

					authArtifacts, ok := authSession.Values["tokens"]
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
					tokenSource := authenticator.GetTokenSource(context.Background(), token)
					newToken, err := tokenSource.Token()
					if err != nil {
						log.Warn().Err(err).Msg("refresh token failed")
					}
					if newToken.AccessToken != token.AccessToken {
						// we need to save this one into the session
						authTokensB, err := json.Marshal(token)
						if err == nil {
							authSession.Values["tokens"] = string(authTokensB)
						}
					}

					break
				}

				// we don't want to create a new session if nobody every created one before
				// we are only here to ensure that the session is an old one and slide it out.
				// i.e. bump out the expiration time
				authSession.Save(c.Request(), c.Response())
			}
			return next(c)
		}
	}
}
