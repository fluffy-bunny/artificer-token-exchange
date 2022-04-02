package session

import (
	"context"
	"echo-starter/internal/session"

	contracts_auth "echo-starter/internal/contracts/auth"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"
	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

const middlewareLogName = "ensure-auth-token-refresh"

func EnsureAuthTokenRefresh(container di.Container) echo.MiddlewareFunc {
	authenticator := core_contracts_oidc.GetIOIDCAuthenticatorFromContainer(container)
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			scopedContainer := c.Get(core_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			logger := contracts_logger.GetILoggerFromContainer(scopedContainer)
			warnEvent := logger.GetLogger().Warn().Str("middleware", middlewareLogName)
			errorEvent := logger.GetLogger().Error().Str("middleware", middlewareLogName)
			debugEvent := logger.GetLogger().Debug().Str("middleware", middlewareLogName)

			for {
				// 1. get our idompontent session
				sess := session.GetSession(c)
				idompotencyKey, ok := sess.Values["idempotency_key"]
				if !ok {
					// if we don't  have this the user hasn't logged in
					break
				}
				tokenStore := contracts_auth.GetIInternalTokenStoreFromContainer(scopedContainer)

				token, err := tokenStore.GetTokenByIdempotencyKey(idompotencyKey.(string))
				if err != nil {
					// not necessarily an error. The tokens could have been removed and our idompotent key could be stale
					debugEvent.Err(err).Msg("Failed to get token")
					break
				}
				tokenSource := authenticator.GetTokenSource(context.Background(), token)
				newToken, err := tokenSource.Token()
				if err != nil {
					warnEvent.Err(err).Msg("refresh token failed")
					break
				}
				if newToken.AccessToken != token.AccessToken {
					err = tokenStore.StoreTokenByIdempotencyKey(idompotencyKey.(string), newToken)
					if err != nil {
						errorEvent.Err(err).Msg("Failed to store token")
					}
				} else {
					err = tokenStore.SlideOutExpiration()
					if err != nil {
						errorEvent.Err(err).Msg("Failed to slide out expiration")
					}
				}
				break
			}

			return next(c)
		}
	}
}
