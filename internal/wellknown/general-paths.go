package wellknown

const (
	HomePath               = "/"
	AboutPath              = "/about"
	HealthzPath            = "/healthz"
	ReadyPath              = "/ready"
	DeepPath               = "/deep/:id/:name"
	ProfilesPath           = "/profiles"
	ArtistsPath            = "/artists"
	APIArtistsPath         = "/api/:version/artists"
	APIArtistsIdPath       = "/api/:version/artists/:id"
	APIArtistsIdAlbumsPath = "/api/:version/artists/:id/albums"
	AccountsPath           = "/accounts"
	APIDevPath             = "/api/:version/dev"

	APIAccountsPath     = "/api/:version/accounts"
	WebHookPath         = "/api/:version/webhook"
	GraphQLEndpointPath = "/api/:version/graphql"
	GraphiQLPath        = "/graphiql"
)
