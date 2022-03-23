package login

import (
	"crypto/rand"
	auth_authenticator "echo-starter/internal/auth/authenticator"
	auth_shared "echo-starter/internal/handlers/auth/shared"
	"echo-starter/internal/session"
	"encoding/base64"
	"encoding/json"
	"net/http"

	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	"github.com/labstack/echo/v4"
)

func Handler(auth *auth_authenticator.Authenticator) func(c echo.Context) error {
	return func(c echo.Context) error {
		u := new(auth_shared.LoginParms)
		if err := c.Bind(u); err != nil {
			return err
		}
		if core_utils.IsEmptyOrNil(u.RedirectURL) {
			u.RedirectURL = "/"
		}
		jsonLoginParams, err := json.Marshal(u)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		state, err := generateRandomState()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())

		}
		sess := session.GetSession(c)
		sess.Values[auth_shared.AuthStateSessionKey] = state
		sess.Values[auth_shared.LoginParamsSessionKey] = jsonLoginParams

		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		url := auth.Config.AuthCodeURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
		return nil
	}
}
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
