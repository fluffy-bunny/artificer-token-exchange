package callback

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	auth_shared "echo-starter/internal/contracts/auth/shared"
	contracts_handler "echo-starter/internal/contracts/handler"
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"encoding/json"
	"net/http"
	"reflect"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Authenticator contracts_auth.IOIDCAuthenticator `inject:"authenticator"`
	}
)

func assertImplementation() {
	var _ contracts_handler.IHandler = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder *di.Builder) {
	contracts_handler.AddScopedIHandlerEx(builder,
		reflectType,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown.OIDCCallbackPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) Do(c echo.Context) error {
	request := c.Request()
	ctx := request.Context()
	state := c.QueryParam("state")
	sess := session.GetSession(c)
	sessionState, _ := sess.Values[auth_shared.AuthStateSessionKey]
	jsonLoginParams, _ := sess.Values[auth_shared.LoginParamsSessionKey]

	loginParams := &auth_shared.LoginParms{}
	json.Unmarshal(jsonLoginParams.([]byte), loginParams)

	if state != sessionState {
		return c.String(http.StatusBadRequest, "Invalid state parameter")
	}

	// Exchange an authorization code for a token.
	token, err := s.Authenticator.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
	}

	idToken, err := s.Authenticator.VerifyIDToken(ctx, token)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to verify ID Token.")

	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	sess.Values["authenticated"] = true

	for k, v := range profile {
		sess.Values["id:"+k] = v
	}

	// NOTE: I have NEVER had the need to store an access token to what is a simple Authentication service. i.e. proof of life.
	// imagine having you website login to google, but you aren't actually using any google services.  The services you are using are yours.
	// So any access_token google gives you during login is only good to get user profile infomration and that is IT.
	// Everything we need we already have in that ID_TOKEN, which we have turned into a stored profile
	// sess.Values["access_token"] = token.AccessToken

	// now that we have logged in we don't need those login paramaters anymore
	delete(sess.Values, auth_shared.LoginParamsSessionKey)
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Redirect to logged in page.
	c.Redirect(http.StatusTemporaryRedirect, loginParams.RedirectURL)
	return nil
}
