package home

import (
	contracts_handler "echo-starter/internal/contracts/handler"

	"echo-starter/internal/wellknown"

	"reflect"
	"strings"

	auth_shared "echo-starter/internal/contracts/auth/shared"
	"echo-starter/internal/models"
	"echo-starter/internal/session"
	"encoding/json"
	"net/http"

	di "github.com/fluffy-bunny/sarulabsdi"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

type (
	service struct {
	}
)

func assertImplementation() {
	var _ contracts_handler.IHandler = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder *di.Builder) {
	httpVerbs := []contracts_handler.HTTPVERB{
		contracts_handler.GET,
	}
	httpVerbS := []string{}
	for _, httpVerb := range httpVerbs {
		httpVerbS = append(httpVerbS, httpVerb.String())
	}
	path := wellknown.HomePath
	metadata := map[string]interface{}{
		"path":      path,
		"httpVerbs": httpVerbs,
	}

	log.Info().
		Str("DI", "IHandler").
		Str("path", path).
		Str("httpVerbs", strings.Join(httpVerbS, "|")).Send()
	contracts_handler.AddScopedIHandlerWithMetadata(builder, reflectType, metadata)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {
	test := &struct {
		Bob string
	}{
		Bob: "bob",
	}
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
			"user":  user,
			"paths": models.NewPaths(),
			"test":  test,
		})
		//	return c.String(http.StatusOK, jsonProfileS)
	}
	return c.Render(http.StatusOK, "content", map[string]interface{}{
		"paths": models.NewPaths(),
		"test":  test,
	})
}
