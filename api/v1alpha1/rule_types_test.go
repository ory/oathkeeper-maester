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
	"github.com/ory/oathkeeper-k8s-controller/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

var (
	template = `[
  {
    "upstream": {
      "url": "http://my-backend-service1",
      "strip_path": "/api/v1",
      "preserve_host": true
    },
    "id": "foo1.default",
    "match": {
      "url": "http://my-app/some-route1",
      "methods": [
        "GET",
        "POST"
      ]
    },
    "authenticators": [
      {
        "handler": "handler1",
        "config": {
          "key1": "val1"
        }
      }
    ],
    "authorizer": {
      "handler": "deny"
    },
    "mutator": {
      "handler": "handler2",
      "config": {
        "key1": [
          "val1",
          "val2",
          "val3"
        ]
      }
    }
  },
  {
    "upstream": {
      "url": "http://my-backend-service2",
      "preserve_host": false
    },
    "id": "foo2.default",
    "match": {
      "url": "http://my-app/some-route2",
      "methods": [
        "GET",
        "POST"
      ]
    },
    "authenticators": [
      {
        "handler": "handler1",
        "config": {
          "key1": "val1"
        }
      },
      {
        "handler": "handler2",
        "config": {
          "key1": [
            "val1",
            "val2",
            "val3"
          ]
        }
      }
    ],
    "authorizer": {
      "handler": "deny"
    },
    "mutator": {
      "handler": "noop"
    }
  },
  {
    "upstream": {
      "url": "http://my-backend-service3",
      "preserve_host": false
    },
    "id": "foo3.default",
    "match": {
      "url": "http://my-app/some-route3",
      "methods": [
        "GET",
        "POST"
      ]
    },
    "authenticators": [
      {
        "handler": "unauthorized"
      }
    ],
    "authorizer": {
      "handler": "handler1",
      "config": {
        "key1": "val1"
      }
    },
    "mutator": {
      "handler": "noop"
    }
  }
]`

	sampleConfig = `{
  "key1": "val1"
}
`

	sampleConfig2 = `{
  "key1": [
    "val1",
    "val2",
    "val3"
  ]
}
`
)

func TestToOathkeeperRules(t *testing.T) {

	t.Run("Should convert a RuleList object into a valid JSON array", func(t *testing.T) {

		var list = &RuleList{}

		t.Run("with no elements if the 'Item' field in the RuleList object is empty", func(t *testing.T) {

			//given
			list.Items = []Rule{}

			//when
			raw, err := list.ToOathkeeperRules()

			//then
			require.NoError(t, err)
			assert.Equal(t, "[]", string(raw))
		})

		t.Run("with raw Oathkeeper rule(s) if the 'Item' field in the RuleList object is not empty", func(t *testing.T) {

			//given
			h1 := newHandler("handler1", sampleConfig)
			h2 := newHandler("handler2", sampleConfig2)

			rule1 := newRule(
				"foo1",
				"default",
				"http://my-backend-service1",
				"http://my-app/some-route1",
				newStringPtr("/api/v1"),
				newBoolPtr(true),
				[]*Authenticator{&Authenticator{h1}},
				nil,
				&Mutator{h2})

			rule2 := newRule(
				"foo2",
				"default",
				"http://my-backend-service2",
				"http://my-app/some-route2",
				nil,
				newBoolPtr(false),
				[]*Authenticator{&Authenticator{h1}, {h2}},
				nil,
				nil)

			rule3 := newRule(
				"foo3",
				"default",
				"http://my-backend-service3",
				"http://my-app/some-route3",
				nil,
				nil,
				nil,
				&Authorizer{h1},
				nil)

			list.Items = []Rule{*rule1, *rule2, *rule3}

			//when
			raw, err := list.ToOathkeeperRules()

			//then
			require.NoError(t, err)
			assert.Equal(t, template, string(raw))
		})
	})
}

func TestToRuleJson(t *testing.T) {

	t.Run("Should convert a Rule to JSON Rule", func(t *testing.T) {

		var actual *RuleJSON
		var testHandler = newHandler("test-handler", "")
		var testRule = newRule(
			"r1",
			"test",
			"https://upstream.url",
			"https://match.this/url",
			newStringPtr("/strip/me"),
			nil,
			nil,
			nil,
			nil)

		t.Run("If no handlers have been specified, it should generate an ID and add default values for missing handlers", func(t *testing.T) {

			//when
			actual = testRule.ToRuleJSON()

			//then
			assert.Equal(t, "r1.test", actual.ID)

			assertHasDefaultAuthenticator(t, actual)

			require.NotNil(t, actual.RuleSpec.Authorizer)
			assert.Equal(t, denyHandler, actual.RuleSpec.Authorizer.Handler)

			require.NotNil(t, actual.RuleSpec.Mutator)
			assert.Equal(t, noopHandler, actual.RuleSpec.Mutator.Handler)

			assert.False(t, *actual.RuleSpec.Upstream.PreserveHost)
		})

		t.Run("If one handler has been provided, it should generate an ID, rewrite the provided handler and add default values for missing handlers", func(t *testing.T) {

			//given
			testRule.Spec.Mutator = &Mutator{testHandler}

			//when
			actual = testRule.ToRuleJSON()

			//then
			assert.Equal(t, "r1.test", actual.ID)

			assertHasDefaultAuthenticator(t, actual)

			require.NotNil(t, actual.RuleSpec.Authorizer)
			assert.Equal(t, denyHandler, actual.RuleSpec.Authorizer.Handler)

			require.NotNil(t, actual.RuleSpec.Mutator)
			assert.Equal(t, testHandler, actual.RuleSpec.Mutator.Handler)
		})

		t.Run("If all handlers are defined, it should generate an ID and rewrite the entire spec", func(t *testing.T) {

			//given
			testRule.Spec.Authenticators = []*Authenticator{{testHandler}}
			testRule.Spec.Authorizer = &Authorizer{testHandler}

			//when
			actual = testRule.ToRuleJSON()

			//then
			assert.Equal(t, "r1.test", actual.ID)
			assert.Equal(t, testRule.Spec, actual.RuleSpec)
		})
	})
}

func TestValidateWith(t *testing.T) {

	var validationError error
	var testHandler = newHandler("handler1", sampleConfig)
	var rule = newRule(
		"foo1",
		"default",
		"http://my-backend-service1",
		"http://my-app/some-route1",
		newStringPtr("/api/v1"),
		newBoolPtr(true),
		nil,
		nil,
		nil)

	var validationConfig = validation.Config{
		AuthenticatorsAvailable: []string{testHandler.Name},
		AuthorizersAvailable:    []string{testHandler.Name},
		MutatorsAvailable:       []string{testHandler.Name},
	}

	t.Run("Should return no error for a rule with", func(t *testing.T) {

		t.Run("no handlers", func(t *testing.T) {

			//when
			validationError = rule.ValidateWith(validationConfig)

			//then
			require.NoError(t, validationError)
		})

		t.Run("allowed handlers", func(t *testing.T) {

			//given
			rule.Spec.Authenticators = []*Authenticator{{testHandler}}
			rule.Spec.Authorizer = &Authorizer{testHandler}
			rule.Spec.Mutator = &Mutator{testHandler}

			//when
			validationError = rule.ValidateWith(validationConfig)

			//then
			require.NoError(t, validationError)
		})
	})

	t.Run("Should return an error for a rule with", func(t *testing.T) {

		t.Run("forbidden handler(s)", func(t *testing.T) {

			//given
			invalidTestHandler := newHandler("notValidHandlerName", sampleConfig)
			rule.Spec.Authenticators = []*Authenticator{{testHandler}, {invalidTestHandler}}

			//when
			validationError = rule.ValidateWith(validationConfig)

			//then
			require.Error(t, validationError)
		})
	})
}

func TestFilterNotValid(t *testing.T) {

	t.Run("Should return only valid rules", func(t *testing.T) {

		//given
		rule1 := newRuleWithStatusOnly(false, newStringPtr("authenticator: sample is invalid"))
		rule2 := newRuleWithStatusOnly(false, nil)
		rule3 := newRuleWithStatusOnly(true, nil)

		list := &RuleList{Items: []Rule{rule1, rule2, rule3}}

		//when
		validationResult := list.FilterNotValid().Items

		//then
		require.NotEmpty(t, validationResult)
		require.Len(t, validationResult, 1)

		assert.Equal(t, rule3, validationResult[0])
	})
}

func newRule(name, namespace, upstreamURL, matchURL string, stripURLPath *string, preserveURLHost *bool, authenticators []*Authenticator, authorizer *Authorizer, mutator *Mutator) *Rule {

	spec := RuleSpec{
		Upstream: &Upstream{
			URL:          upstreamURL,
			PreserveHost: preserveURLHost,
			StripPath:    stripURLPath,
		},
		Match: &Match{
			URL:     matchURL,
			Methods: []string{"GET", "POST"},
		},
		Authenticators: authenticators,
		Authorizer:     authorizer,
		Mutator:        mutator,
	}

	return &Rule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: spec,
	}
}

func newRuleWithStatusOnly(valid bool, validationError *string) Rule {
	return Rule{
		Status: RuleStatus{
			Validation: &Validation{
				Valid: &valid,
				Error: validationError,
			},
		},
	}
}

func newHandler(name, config string) *Handler {
	h := &Handler{
		Name: name,
	}

	if config != "" {
		h.Config = &runtime.RawExtension{
			Raw: []byte(config),
		}
	}

	return h
}

func newBoolPtr(b bool) *bool {
	return &b
}

func newStringPtr(s string) *string {
	return &s
}

func assertHasDefaultAuthenticator(t *testing.T, actual *RuleJSON) {
	require.NotNil(t, actual.RuleSpec.Authenticators)
	require.NotEmpty(t, actual.RuleSpec.Authenticators)
	require.Len(t, actual.RuleSpec.Authenticators, 1)
	assert.Equal(t, unauthorizedHandler, actual.RuleSpec.Authenticators[0].Handler)
}
