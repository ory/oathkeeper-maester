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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ory/oathkeeper-maester/internal/validation"
)

const (
	deny         = "deny"
	noop         = "noop"
	unauthorized = "unauthorized"
)

var (
	denyHandler         = &Handler{Name: deny}
	noopHandler         = &Handler{Name: noop}
	unauthorizedHandler = &Handler{Name: unauthorized}
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
	// +kubebuilder:validation:Optional
	// +optional
	Upstream       *Upstream        `json:"upstream,omitempty"`
	Match          *Match           `json:"match"`
	Authenticators []*Authenticator `json:"authenticators,omitempty"`
	Authorizer     *Authorizer      `json:"authorizer,omitempty"`
	Mutators       []*Mutator       `json:"mutators,omitempty"`
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*
	//
	// ConfigMapName points to the K8s ConfigMap that contains these rules
	ConfigMapName *string `json:"configMapName,omitempty"`
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
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Pattern=`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`
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
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:XPreserveUnknownFields
	Config *runtime.RawExtension `json:"config,omitempty"`
}

// ToOathkeeperRules transforms a RuleList object into a JSON object digestible by Oathkeeper.
func (rl RuleList) ToOathkeeperRules() ([]byte, error) {

	rules := make([]*RuleJSON, len(rl.Items))

	for i := range rl.Items {
		rules[i] = rl.Items[i].ToRuleJSON()
	}

	return unescapedMarshalIndent(rules, "", "  ")
}

// FilterNotValid filters out Rules which doesn't pass validation due to being not processed yet or due to negative result of validation. It returns a list of Rules which passed validation successfully.
func (rl RuleList) FilterNotValid() RuleList {
	rlCopy := rl
	validRules := []Rule{}
	for _, rule := range rl.Items {
		if rule.Status.Validation != nil && rule.Status.Validation.Valid != nil && *rule.Status.Validation.Valid {
			validRules = append(validRules, rule)
		}
	}
	rlCopy.Items = validRules
	return rlCopy
}

// FilterConfigMapName filters out Rules that don't effect the given ConfigMap
func (rl RuleList) FilterConfigMapName(name *string) RuleList {
	rlCopy := rl
	validRules := []Rule{}
	for _, rule := range rl.Items {
		if rule.Spec.ConfigMapName == nil {
			if name == nil {
				validRules = append(validRules, rule)
			}
		} else if *rule.Spec.ConfigMapName == *name {
			validRules = append(validRules, rule)
		}
	}
	rlCopy.Items = validRules
	return rlCopy
}

// FilterOutRule filters out the provided rule from the rule list, for re-generating the rules when a rule is deleted
func (rl RuleList) FilterOutRule(r Rule) RuleList {
	rlCopy := rl
	validRules := []Rule{}
	for _, rule := range rl.Items {
		if rule.ObjectMeta.UID != r.ObjectMeta.UID {
			validRules = append(validRules, rule)
		}
	}
	rlCopy.Items = validRules
	return rlCopy
}

// ValidateWith uses provided validation configuration to check whether the rule have proper handlers set. Nil is a valid handler.
func (r Rule) ValidateWith(config validation.Config) error {

	var invalidHandlers []string

	if r.Spec.Authenticators != nil {
		for _, authenticator := range r.Spec.Authenticators {
			if valid := config.IsAuthenticatorValid(authenticator.Name); !valid {
				invalidHandlers = append(invalidHandlers, fmt.Sprintf("authenticator/%s", authenticator.Name))
			}
		}
	}

	if r.Spec.Authorizer != nil {
		if valid := config.IsAuthorizerValid(r.Spec.Authorizer.Name); !valid {
			invalidHandlers = append(invalidHandlers, fmt.Sprintf("authorizer/%s", r.Spec.Authorizer.Name))
		}
	}

	if r.Spec.Mutators != nil {
		for _, m := range r.Spec.Mutators {
			if valid := config.IsMutatorValid(m.Name); !valid {
				invalidHandlers = append(invalidHandlers, m.Name)
			}
		}
	}

	if len(invalidHandlers) != 0 {
		return fmt.Errorf("invalid handlers: %s, please check the configuration", invalidHandlers)
	}

	return nil
}

// ToRuleJSON transforms a Rule object into an intermediary RuleJSON object
func (r Rule) ToRuleJSON() *RuleJSON {

	ruleJSON := &RuleJSON{
		ID:       r.Name + "." + r.Namespace,
		RuleSpec: r.Spec,
	}

	if ruleJSON.Authenticators == nil {
		ruleJSON.Authenticators = []*Authenticator{{unauthorizedHandler}}
	}
	if ruleJSON.Authorizer == nil {
		ruleJSON.Authorizer = &Authorizer{denyHandler}
	}
	if ruleJSON.Mutators == nil {
		ruleJSON.Mutators = []*Mutator{{noopHandler}}
	}

	if ruleJSON.Upstream == nil {
		ruleJSON.Upstream = &Upstream{}
	}

	if ruleJSON.Upstream.PreserveHost == nil {
		ruleJSON.Upstream.PreserveHost = &preserveHostDefault
	}

	return ruleJSON
}

func init() {
	SchemeBuilder.Register(&Rule{}, &RuleList{})
}
