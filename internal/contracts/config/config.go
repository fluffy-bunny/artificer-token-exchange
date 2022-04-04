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
	// Config type
	Config struct {
		ApplicationName           string     `json:"applicationName" mapstructure:"APPLICATION_NAME"`
		ApplicationEnvironment    string     `json:"applicationEnvironment" mapstructure:"APPLICATION_ENVIRONMENT"`
		PrettyLog                 bool       `json:"prettyLog" mapstructure:"PRETTY_LOG"`
		LogLevel                  string     `json:"logLevel" mapstructure:"LOG_LEVEL"`
		Port                      int        `json:"port" mapstructure:"PORT"`
		Oidc                      oidcConfig `json:"oidc" mapstructure:"OIDC"`
		SessionMaxAgeSeconds      int        `json:"sessionMaxAgeSeconds" mapstructure:"SESSION_MAX_AGE_SECONDS"`
		AuthCookieExpireSeconds   int        `json:"authCookieExpireSeconds" mapstructure:"AUTH_COOKIE_EXPIRE_SECONDS"`
		AuthCookieName            string     `json:"authCookieName" mapstructure:"AUTH_COOKIE_NAME"`
		SecureCookieHashKey       string     `json:"secureCookieHashKey" mapstructure:"SECURE_COOKIE_HASH_KEY"`
		SecureCookieEncryptionKey string     `json:"secureCookieEncryptionKey" mapstructure:"SECURE_COOKIE_ENCRYPTION_KEY"`
		GraphQLEndpoint           string     `json:"graphQLEndpoint" mapstructure:"GRAPHQL_ENDPOINT"`
		// cookie|inmemory|redis
		SessionEngine string `json:"sessionEngine" mapstructure:"SESSION_ENGINE"`
		RedisUrl      string `json:"redisUrl" mapstructure:"REDIS_URL"`
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
 	"SESSION_MAX_AGE_SECONDS": 60,
	"AUTH_COOKIE_EXPIRE_SECONDS": 60,
	"AUTH_COOKIE_NAME": "_auth",
	"SECURE_COOKIE_HASH_KEY": "",
	"SECURE_COOKIE_ENCRYPTION_KEY": "",
	"GRAPHQL_ENDPOINT": "https://countries.trevorblades.com/",
	"SESSION_ENGINE": "cookie",
	"REDIS_URL": "redis://localhost:6379"


}
`)
)
