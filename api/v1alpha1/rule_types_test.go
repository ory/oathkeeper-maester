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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// These tests are written in BDD-style using Ginkgo framework. Refer to
// http://onsi.github.io/ginkgo to learn more.

var _ = Describe("Rule", func() {
	var (
		key              types.NamespacedName
		created, fetched *Rule
	)

	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	// Add Tests for OpenAPI validation (or additonal CRD features) specified in
	// your API definition.
	// Avoid adding tests for vanilla CRUD operations because they would
	// test Kubernetes API server, which isn't the goal here.
	Context("Create API", func() {

		It("should create an object successfully", func() {

			key = types.NamespacedName{
				Name:      "foo",
				Namespace: "default",
			}

			t := true

			h := &Handler{
				Name: "sample-handler",
				Config: &runtime.RawExtension{
					Raw: []byte("{}"),
				},
			}

			rs := RuleSpec{
				ID: "sample-rule1",
				Upstream: &Upstream{
					URL:          "https://url.com",
					PreserveHost: &t,
				},
				Match: &Match{
					URL:     "https://url2.com",
					Methods: []string{"GET", "POST"},
				},
				Authenticators: []*Authenticator{&Authenticator{h}},
				Authorizer:     &Authorizer{h},
				Mutator:        &Mutator{h},
			}

			created = &Rule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
				},
				Spec:   rs,
				Status: RuleStatus{},
			}

			By("creating an API obj")
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			fetched = &Rule{}
			Expect(k8sClient.Get(context.TODO(), key, fetched)).To(Succeed())
			Expect(fetched).To(Equal(created))

			By("deleting the created object")
			Expect(k8sClient.Delete(context.TODO(), created)).To(Succeed())
			Expect(k8sClient.Get(context.TODO(), key, created)).ToNot(Succeed())
		})
	})

	var template = `[
  {
    "id": "some-id-1",
    "upstream": {
      "url": "http://my-backend-service1",
      "strip_path": "/api/v1",
      "preserve_host": true
    },
    "match": {
      "url": "http://my-app/some-route1",
      "methods": [
        "GET",
        "POST"
      ]
    },
    "authenticators": [
      {
        "handler": "allow"
      }
    ],
    "mutator": {
      "handler": "noop"
    }
  },
  {
    "id": "some-id-2",
    "upstream": {
      "url": "http://my-backend-service2",
      "preserve_host": false
    },
    "match": {
      "url": "http://my-app/some-route2",
      "methods": [
        "GET",
        "POST"
      ]
    },
    "authenticators": [
      {
        "handler": "allow"
      },
      {
        "handler": "noop"
      }
    ]
  }
]`

	Context("ToOathkeeperRules", func() {

		It("Should return a JSON array of raw Oathkeeper rules", func() {

			s := "/api/v1"

			t1 := true
			t2 := false

			h1 := &Handler{
				Name: "allow",
			}

			h2 := &Handler{
				Name: "noop",
			}

			rs1 := RuleSpec{
				ID: "some-id-1",
				Upstream: &Upstream{
					URL:          "http://my-backend-service1",
					StripPath:    &s,
					PreserveHost: &t1,
				},
				Match: &Match{
					URL:     "http://my-app/some-route1",
					Methods: []string{"GET", "POST"},
				},
				Authenticators: []*Authenticator{&Authenticator{h1}},
				Mutator:        &Mutator{h2},
			}

			rs2 := RuleSpec{
				ID: "some-id-2",
				Upstream: &Upstream{
					URL:          "http://my-backend-service2",
					PreserveHost: &t2,
				},
				Match: &Match{
					URL:     "http://my-app/some-route2",
					Methods: []string{"GET", "POST"},
				},
				Authenticators: []*Authenticator{&Authenticator{h1}, {h2}},
			}

			created1 := Rule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo1",
					Namespace: "default",
				},
				Spec:   rs1,
				Status: RuleStatus{},
			}

			created2 := Rule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo2",
					Namespace: "default",
				},
				Spec:   rs2,
				Status: RuleStatus{},
			}

			list := &RuleList{Items: []Rule{created1, created2}}

			By("transforming the receiver into a slice of bytes")

			raw, err := list.ToOathkeeperRules()

			Expect(err).To(BeNil())
			Expect(string(raw)).To(Equal(template))

		})
	})
})
