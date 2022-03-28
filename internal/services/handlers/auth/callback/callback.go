package callback

import (
	contracts_auth "echo-starter/internal/contracts/auth"
	auth_shared "echo-starter/internal/contracts/auth/shared"
	contracts_claimsprovider "echo-starter/internal/contracts/claimsprovider"
	contracts_handler "echo-starter/internal/contracts/handler"
	"echo-starter/internal/session"
	"echo-starter/internal/wellknown"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	contracts_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Authenticator  contracts_auth.IOIDCAuthenticator        `inject:""`
		AuthCookie     contracts_auth.IAuthCookie               `inject:""`
		ClaimsProvider contracts_claimsprovider.IClaimsProvider `inject:""`
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
	jsonTokenB, err := json.Marshal(token)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	sess.Values["_authArtifacts"] = string(jsonTokenB)

	idToken, err := s.Authenticator.VerifyIDToken(ctx, token)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to verify ID Token.")
	}

	profileClaims, err := s.ClaimsProvider.GetClaims(idToken.Subject, "")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch claims.")
	}
	jsonBytes, err := json.Marshal(profileClaims)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to marshal profileClaims.")
	}
	sess.Values["_profile"] = string(jsonBytes)

	var identityProfile map[string]interface{}
	if err := idToken.Claims(&identityProfile); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	var identityClaims []*contracts_claimsprincipal.Claim
	for k, v := range identityProfile {
		switch typed := v.(type) {
		case string:
			identityClaims = append(identityClaims, &contracts_claimsprincipal.Claim{
				Type:  k,
				Value: typed,
			})
		case float64:
			identityClaims = append(identityClaims, &contracts_claimsprincipal.Claim{
				Type:  k,
				Value: fmt.Sprintf("%f", typed),
			})

		}
	}

	// we store the claims in the session in context of the userid
	jsonBytes, err = json.Marshal(identityClaims)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to marshal identityClaims.")
	}
	sess.Values[idToken.Subject] = string(jsonBytes)

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
	// our auth cookie simply stores the userid which points to the entry in the session
	// this is to prepare for when the session is backed by a session backend store and not a fat cookie store
	s.AuthCookie.SetAuthCookieValue(c, idToken.Subject)

	// Redirect to logged in page.
	c.Redirect(http.StatusFound, loginParams.RedirectURL)
	return nil
}
