package user

import (
	auth_authenticator "echo-starter/internal/auth/authenticator"
	auth_shared "echo-starter/internal/handlers/auth/shared"
	"echo-starter/internal/session"
	"echo-starter/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler(auth *auth_authenticator.Authenticator) func(c echo.Context) error {
	return func(c echo.Context) error {
		sess := session.GetSession(c)
		jsonProfile, _ := sess.Values[auth_shared.ProfileSessionKey]
		if jsonProfile != nil {
			var profile map[string]interface{}
			json.Unmarshal(jsonProfile.([]byte), &profile)
			jsonProfileS := utils.PrettyJSON(profile)
			return c.String(http.StatusOK, jsonProfileS)
		} else {
			return c.String(http.StatusOK, "No profile found")
		}
	}
}
