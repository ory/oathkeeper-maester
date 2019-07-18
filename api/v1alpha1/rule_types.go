/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"encoding/json"
	"fmt"
	"github.com/ory/oathkeeper-k8s-controller/internal/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	allow = "allow"
	noop  = "noop"
)

var (
	noopHandler         = &Handler{Name: noop}
	allowHandler        = &Handler{Name: allow}
	preserveHostDefault = false
)

// +kubebuilder:object:root=true
// Rule is the Schema for the rules API
type Rule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RuleSpec   `json:"spec,omitempty"`
	Status RuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// RuleList contains a list of Rule
type RuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rule `json:"items"`
}

// RuleSpec defines the desired state of Rule
type RuleSpec struct {
	Upstream       *Upstream        `json:"upstream"`
	Match          *Match           `json:"match"`
	Authenticators []*Authenticator `json:"authenticators,omitempty"`
	Authorizer     *Authorizer      `json:"authorizer,omitempty"`
	Mutator        *Mutator         `json:"mutator,omitempty"`
}

// Validation defines the validation state of Rule
type Validation struct {
	// +optional
	Valid *bool `json:"valid,omitempty"`
	// +optional
	Error *string `json:"validationError,omitempty"`
}

// RuleStatus defines the observed state of Rule
type RuleStatus struct {
	// +optional
	Validation *Validation `json:"validation,omitempty"`
}

// Upstream represents the location of a server where requests matching a rule should be forwarded to.
type Upstream struct {
	// URL defines the target URL for incoming requests
	URL string `json:"url"`
	// StripPath replaces the provided path prefix when forwarding the requested URL to the upstream URL.
	// +optional
	StripPath *string `json:"stripPath,omitempty"`
	// PreserveHost includes the host and port of the url value if set to false. If true, the host and port of the ORY Oathkeeper Proxy will be used instead.
	// +optional
	PreserveHost *bool `json:"preserveHost,omitempty"`
}

// Match defines the URL(s) that an access rule should match.
type Match struct {
	// URL is the URL that should be matched. It supports regex templates.
	URL string `json:"url"`
	// Methods represent an array of HTTP methods (e.g. GET, POST, PUT, DELETE, ...)
	Methods []string `json:"methods"`
}

// Authenticator represents a handler that authenticates provided credentials.
type Authenticator struct {
	*Handler `json:",inline"`
}

// Authorizer represents a handler that authorizes the subject ("user") from the previously validated credentials making the request.
type Authorizer struct {
	*Handler `json:",inline"`
}

// Mutator represents a handler that transforms the HTTP request before forwarding it.
type Mutator struct {
	*Handler `json:",inline"`
}

// Handler represents an Oathkeeper routine that operates on incoming requests. It is used to either validate a request (Authenticator, Authorizer) or modify it (Mutator).
type Handler struct {
	// Name is the name of a handler
	Name string `json:"handler"`
	// Config configures the handler. Configuration keys vary per handler.
	// +kubebuilder:validation:Type=object
	Config *runtime.RawExtension `json:"config,omitempty"`
}

// ToOathkeeperRules transforms a RuleList object into a JSON object digestible by Oathkeeper
func (rl RuleList) ToOathkeeperRules() ([]byte, error) {

	var rules []*RuleJSON

	for _, item := range rl.Items {
		rules = append(rules, item.ToRuleJSON())
	}

	return json.MarshalIndent(rules, "", "  ")
}

// ToRuleJSON transforms a Rule object into an intermediary RuleJSON object
func (r Rule) ToRuleJSON() *RuleJSON {

	ruleJSON := &RuleJSON{
		ID:       r.Name + "." + r.Namespace,
		RuleSpec: r.Spec,
	}

	if ruleJSON.Authenticators == nil {
		ruleJSON.Authenticators = []*Authenticator{{noopHandler}}
	}
	if ruleJSON.Authorizer == nil {
		ruleJSON.Authorizer = &Authorizer{allowHandler}
	}
	if ruleJSON.Mutator == nil {
		ruleJSON.Mutator = &Mutator{noopHandler}
	}

	if ruleJSON.Upstream.PreserveHost == nil {
		ruleJSON.Upstream.PreserveHost = &preserveHostDefault
	}

	return ruleJSON
}

func init() {
	SchemeBuilder.Register(&Rule{}, &RuleList{})
}
