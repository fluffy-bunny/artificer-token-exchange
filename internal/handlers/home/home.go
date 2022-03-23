package home

import (
	auth_shared "echo-starter/internal/handlers/auth/shared"
	"echo-starter/internal/models"
	"echo-starter/internal/session"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler() func(c echo.Context) error {
	return func(c echo.Context) error {
		sess := session.GetSession(c)
		jsonProfile, _ := sess.Values[auth_shared.ProfileSessionKey]
		if jsonProfile != nil {
			var profile map[string]interface{}
			json.Unmarshal(jsonProfile.([]byte), &profile)
			//	jsonProfileS := utils.PrettyJSON(profile)
			user := &models.User{
				Name: profile["name"].(string),
				ID:   profile["sub"].(string),
			}
			return c.Render(http.StatusOK, "content", map[string]interface{}{
				"user": user,
			})
			//	return c.String(http.StatusOK, jsonProfileS)
		}
		return c.Render(http.StatusOK, "content", map[string]interface{}{})

	}
}
