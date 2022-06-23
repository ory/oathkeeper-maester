package v1alpha1

var (
	DefaultAuthenticatorsAvailable = [...]string{"noop", "unauthorized", "anonymous", "cookie_session", "oauth2_client_credentials", "oauth2_introspection", "jwt", "bearer_token"}
	DefaultAuthorizersAvailable    = [...]string{"allow", "deny", "keto_engine_acp_ory", "remote", "remote_json"}
	DefaultMutatorsAvailable       = [...]string{"noop", "id_token", "header", "cookie", "hydrator"}
)

const (
	AuthenticatorsAvailableEnv = "authenticatorsAvailable"
	AuthorizersAvailableEnv    = "authorizersAvailable"
	MutatorsAvailableEnv       = "mutatorsAvailable"
	RulesFileNameRegexp        = "\\A[-._a-zA-Z0-9]+\\z"
)
