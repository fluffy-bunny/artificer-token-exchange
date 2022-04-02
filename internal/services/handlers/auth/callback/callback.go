package callback

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	auth_shared "echo-starter/internal/contracts/auth/shared"
	contracts_claimsprovider "echo-starter/internal/contracts/claimsprovider"
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"encoding/json"
	"net/http"
	"reflect"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger         contracts_logger.ILogger                 `inject:""`
		Authenticator  core_contracts_oidc.IOIDCAuthenticator   `inject:""`
		ClaimsProvider contracts_claimsprovider.IClaimsProvider `inject:""`
		TokenStore     contracts_auth.IInternalTokenStore       `inject:""`
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
	sessionState, _ := sess.Values[auth_shared.AuthStateSessionKey].(string)
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
	_, err = s.Authenticator.VerifyIDToken(ctx, token)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to verify ID Token.")
	}

	/*
		authTokensB, err := json.Marshal(token)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		authSess.Values["tokens"] = string(authTokensB)
	*/
	// the token store and the session share a idompotency key to bind them together
	s.TokenStore.StoreTokenByIdompotencyKey(sessionState, token)
	sess.Values["idompontency_key"] = sessionState

	tt, err := s.TokenStore.GetTokenByIdompotencyKey(sessionState)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if tt.AccessToken != token.AccessToken {
		return c.String(http.StatusInternalServerError, "Token mismatch")
	}
	// now that we have logged in we don't need those login paramaters anymore
	delete(sess.Values, auth_shared.AuthStateSessionKey)
	delete(sess.Values, auth_shared.LoginParamsSessionKey)
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// our auth cookie simply stores the userid which points to the entry in the session
	// this is to prepare for when the session is backed by a session backend store and not a fat cookie store
	//s.AuthCookie.SetAuthCookieValue(c, idToken.Subject)

	// Redirect to logged in page.
	c.Redirect(http.StatusFound, loginParams.RedirectURL)
	return nil
}
