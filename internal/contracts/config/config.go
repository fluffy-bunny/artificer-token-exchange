package config

const (
	Environment_Development = "Development"
)

type (
	oidcConfig struct {
		Domain       string `json:"domain" mapstructure:"DOMAIN"`
		ClientID     string `json:"client_id" mapstructure:"CLIENT_ID"`
		ClientSecret string `json:"client_secret" mapstructure:"CLIENT_SECRET"`
		CallbackURL  string `json:"callback_url" mapstructure:"CALLBACK_URL"`
	}
	oauth2Config struct {
		// ClientID is the application's ID.
		ClientID string `json:"client_id" mapstructure:"CLIENT_ID"`

		// ClientSecret is the application's secret.
		ClientSecret string `json:"client_secret" mapstructure:"CLIENT_SECRET"`

		// RedirectURL is the URL to redirect users going through
		// the OAuth flow, after the resource owner's URLs.
		RedirectURL string `json:"redirect_url" mapstructure:"REDIRECT_URL"`

		// Scope specifies optional requested permissions.
		Scopes []string `json:"scopes" mapstructure:"SCOPES"`
	}
	// Config type
	Config struct {
		ApplicationName         string       `json:"applicationName" mapstructure:"APPLICATION_NAME"`
		ApplicationEnvironment  string       `json:"applicationEnvironment" mapstructure:"APPLICATION_ENVIRONMENT"`
		PrettyLog               bool         `json:"prettyLog" mapstructure:"PRETTY_LOG"`
		LogLevel                string       `json:"logLevel" mapstructure:"LOG_LEVEL"`
		Port                    int          `json:"port" mapstructure:"PORT"`
		OIDC                    oidcConfig   `json:"oidc" mapstructure:"OIDC"`
		OAuth2                  oauth2Config `json:"oauth2" mapstructure:"OAUTH2"`
		SessionMaxAgeSeconds    int          `json:"sessionMaxAgeSeconds" mapstructure:"SESSION_MAX_AGE_SECONDS"`
		AuthCookieExpireSeconds int          `json:"authCookieExpireSeconds" mapstructure:"AUTH_COOKIE_EXPIRE_SECONDS"`
		AuthCookieName          string       `json:"authCookieName" mapstructure:"AUTH_COOKIE_NAME"`
		// session|cookie
		AuthStore                 string `json:"authStore" mapstructure:"AUTH_STORE"`
		SecureCookieHashKey       string `json:"secureCookieHashKey" mapstructure:"SECURE_COOKIE_HASH_KEY"`
		SecureCookieEncryptionKey string `json:"secureCookieEncryptionKey" mapstructure:"SECURE_COOKIE_ENCRYPTION_KEY"`
		GraphQLEndpoint           string `json:"graphQLEndpoint" mapstructure:"GRAPHQL_ENDPOINT"`
		// cookie|inmemory|redis
		SessionEngine string `json:"sessionEngine" mapstructure:"SESSION_ENGINE"`
		RedisUrl      string `json:"redisUrl" mapstructure:"REDIS_URL"`
		RedisPassword string `json:"redisPassword" mapstructure:"REDIS_PASSWORD"`

		// github,oidc
		AuthProvider string `json:"authProvider" mapstructure:"AUTH_PROVIDER"`
	}
)

var (
	// ConfigDefaultJSON default json
	ConfigDefaultJSON = []byte(`
{
	"APPLICATION_NAME": "in-environment",
	"APPLICATION_ENVIRONMENT": "in-environment",
	"PRETTY_LOG": false,
	"LOG_LEVEL": "info",
	"PORT": 1111,
	"OIDC": {
		"DOMAIN": "blah.auth0.com",
		"CLIENT_ID": "in-environment",
		"CLIENT_SECRET": "in-environment",
		"CALLBACK_URL": ""
	},
	"OAUTH2": {
		"CLIENT_ID": "in-environment",
		"CLIENT_SECRET": "in-environment",
		"REDIRECT_URL": "",
		"SCOPES": ""
	},
 	"SESSION_MAX_AGE_SECONDS": 60,
    "AUTH_PROVIDER": "oidc",
	"AUTH_COOKIE_EXPIRE_SECONDS": 60,
	"AUTH_COOKIE_NAME": "_auth",
	"AUTH_STORE": "cookie",
	"SECURE_COOKIE_HASH_KEY": "",
	"SECURE_COOKIE_ENCRYPTION_KEY": "",
	"GRAPHQL_ENDPOINT": "https://countries.trevorblades.com/",
	"SESSION_ENGINE": "cookie",
	"REDIS_URL": "localhost:6379",
	"REDIS_PASSWORD": ""


}
`)
)
