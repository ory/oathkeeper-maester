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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	ID             string           `json:"id"`
	Upstream       *Upstream        `json:"upstream"`
	Match          *Match           `json:"match"`
	Authenticators []*Authenticator `json:"authenticators,omitempty"`
	Authorizer     *Authorizer      `json:"authorizer,omitempty"`
	Mutator        *Mutator         `json:"mutator,omitempty"`
}

// RuleStatus defines the observed state of Rule
type RuleStatus struct {
}

type Upstream struct {
	URL string `json:"url"`
	// +optional
	StripPath *string `json:"stripPath,omitempty"`
	// +optional
	PreserveHost *bool `json:"preserveHost,omitempty"`
}

type Match struct {
	URL     string   `json:"url"`
	Methods []string `json:"methods"`
}

type Authenticator struct {
	*Handler `json:",inline"`
}

type Authorizer struct {
	*Handler `json:",inline"`
}

type Mutator struct {
	*Handler `json:",inline"`
}

type Handler struct {
	Name string `json:"handler"`
	// +kubebuilder:validation:Type=object
	Config *runtime.RawExtension `json:"config,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Rule{}, &RuleList{})
}
