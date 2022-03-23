package config

type oidcConfig struct {
	Domain       string `json:"domain" mapstructure:"DOMAIN"`
	ClientID     string `json:"client_id" mapstructure:"CLIENT_ID"`
	ClientSecret string `json:"client_secret" mapstructure:"CLIENT_SECRET"`
}

// Config type
type Config struct {
	ApplicationName        string     `json:"applicationName" mapstructure:"APPLICATION_NAME"`
	ApplicationEnvironment string     `json:"applicationEnvironment" mapstructure:"APPLICATION_ENVIRONMENT"`
	PrettyLog              bool       `json:"prettyLog" mapstructure:"PRETTY_LOG"`
	LogLevel               string     `json:"logLevel" mapstructure:"LOG_LEVEL"`
	Port                   int        `json:"port" mapstructure:"PORT"`
	Oidc                   oidcConfig `json:"oidc" mapstructure:"OIDC"`
}

// ConfigDefaultJSON default json
var ConfigDefaultJSON = []byte(`
{
	"APPLICATION_NAME": "in-environment",
	"APPLICATION_ENVIRONMENT": "in-environment",
	"PRETTY_LOG": false,
	"LOG_LEVEL": "info",
	"PORT": 1111,
	"OIDC": {
		"DOMAIN": "blah.auth0.com",
		"CLIENT_ID": "in-environment",
		"CLIENT_SECRET": "in-environment"
	}
}
`)
