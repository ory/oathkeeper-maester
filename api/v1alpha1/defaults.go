// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

var (
	DefaultAuthenticatorsAvailable = [...]string{"noop", "unauthorized", "anonymous", "cookie_session", "oauth2_client_credentials", "oauth2_introspection", "jwt", "bearer_token"}
	DefaultAuthorizersAvailable    = [...]string{"allow", "deny", "keto_engine_acp_ory", "remote", "remote_json"}
	DefaultMutatorsAvailable       = [...]string{"noop", "id_token", "header", "cookie", "hydrator"}
	DefaultErrorsAvailable         = [...]string{"json", "redirect", "www_authenticate"}
)

const (
	AuthenticatorsAvailableEnv = "authenticatorsAvailable"
	AuthorizersAvailableEnv    = "authorizersAvailable"
	MutatorsAvailableEnv       = "mutatorsAvailable"
	ErrorsAvailableEnv         = "errorsAvailable"
	RulesFileNameRegexp        = "\\A[-._a-zA-Z0-9]+\\z"
)
