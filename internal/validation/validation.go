// Copyright © 2022 Ory Corp

package validation

type Config struct {
	AuthenticatorsAvailable []string
	AuthorizersAvailable    []string
	MutatorsAvailable       []string
}

func (c Config) IsAuthenticatorValid(authenticator string) bool {
	return isValid(authenticator, c.AuthenticatorsAvailable)
}
func (c Config) IsAuthorizerValid(authorizer string) bool {
	return isValid(authorizer, c.AuthorizersAvailable)
}
func (c Config) IsMutatorValid(mutator string) bool {
	return isValid(mutator, c.MutatorsAvailable)
}

func isValid(current string, available []string) bool {
	for _, a := range available {
		if current == a {
			return true
		}
	}
	return false
}
