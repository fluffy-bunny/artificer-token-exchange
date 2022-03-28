package session

import (
	"echo-starter/internal/session"
	"encoding/json"
	"fmt"
	"time"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func EnsureSlidingSession(container di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			sess := session.GetSession(c)
			if !sess.IsNew {

				authArtifacts, ok := sess.Values["_authArtifacts"]
				if ok {
					token := &oauth2.Token{}
					json.Unmarshal([]byte(authArtifacts.(string)), token)
					if time.Now().Add(5 * time.Minute).After(token.Expiry) {
						// call refresh_token endpoint.
						fmt.Println("refresh token - time")
					}
				}

				// we don't want to create a new session if nobody every created one before
				// we are only here to ensure that the session is an old one and slide it out.
				// i.e. bump out the expiration time
				sess.Save(c.Request(), c.Response())
			}
			return next(c)
		}
	}
}
