// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

// RuleJson is a representation of an Oathkeeper rule.
type RuleJSON struct {
	ID       string `json:"id"`
	RuleSpec `json:",inline"`
}

// MarshalJSON is a custom marshal function that converts RuleJSON objects into JSON objects digestible by Oathkeeper
func (rj RuleJSON) MarshalJSON() ([]byte, error) {

	type Alias RuleJSON

	return unescapedMarshal(&struct {
		Upstream *UpstreamJSON `json:"upstream,omitempty"`
		Alias
	}{
		Upstream: &UpstreamJSON{
			URL:          rj.Upstream.URL,
			PreserveHost: rj.Upstream.PreserveHost,
			StripPath:    rj.Upstream.StripPath,
		},
		Alias: (Alias)(rj),
	})
}

// UpstreamJSON is a helper struct that representats Oathkeeper's upstream object.
type UpstreamJSON struct {
	URL          string  `json:"url"`
	StripPath    *string `json:"strip_path,omitempty"`
	PreserveHost *bool   `json:"preserve_host"`
}
