package authenticator

import (
	"net/http"

	"echo-starter/internal/session"

	"github.com/labstack/echo/v4"
)

func OIDCCallBackHandler(auth *Authenticator) func(c echo.Context) error {
	return func(c echo.Context) error {
		request := c.Request()
		ctx := request.Context()
		state := c.QueryParam("state")
		sess := session.GetSession(c)
		sessionState, _ := sess.Values["state"]

		if state != sessionState {
			return c.String(http.StatusBadRequest, "Invalid state parameter")
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx, c.QueryParam("code"))
		if err != nil {
			return c.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
		}

		idToken, err := auth.VerifyIDToken(ctx, token)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to verify ID Token.")

		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		sess.Values["profile"] = profile
		sess.Values["access_token"] = token.AccessToken

		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Redirect to logged in page.
		c.Redirect(http.StatusTemporaryRedirect, "/user")
		return nil
	}
}
