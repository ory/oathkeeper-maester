package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

type Handler struct {
	Name string `json:"handler"`
	// +kubebuilder:validation:Type=object
	Config *runtime.RawExtension `json:"config,omitempty"`
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
