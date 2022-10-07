// Copyright © 2022 Ory Corp

package integration

import (
	"fmt"

	json "github.com/bitly/go-simplejson"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Validates rules are equal
// `actual` is a representation of an entry from the ConfigMap handled by the Controller
func validateRuleEquals(actual *json.Json, expected *json.Json) {
	Expect(actual).To(Equal(expected))
	expectOnlyKeys(actual, "id", "upstream", "match", "authenticators", "authorizer", "mutators")

	expectString(actual, "id")
	compareUpstreams(actual.Get("upstream"), expected.Get("upstream"))
	compareMatches(actual.Get("match"), expected.Get("match"))
	compareHandlerArrays(actual.Get("authenticators"), expected.Get("authenticators"))
	compareHandlers(actual.Get("authorizer"), expected.Get("authorizer"))
	compareHandlerArrays(actual.Get("mutators"), expected.Get("mutators"))
}

func compareUpstreams(actual *json.Json, expected *json.Json) {
	Expect(actual).To(Equal(expected))
	expectOnlyKeys(actual, "url", "preserve_host")
	expectString(actual, "url")
	expectBoolean(actual.Get("preserve_host"))
}

func compareMatches(actual *json.Json, expected *json.Json) {
	Expect(actual).To(Equal(expected))
	expectOnlyKeys(actual, "url", "methods")
	expectString(actual, "url")
	expectStringArray(actual, "methods")
}

func compareHandlerArrays(actual *json.Json, expected *json.Json) {
	//both are equal
	Expect(actual).To(Equal(expected))

	//expected is an Array
	expectedArray, err := expected.Array()
	Expect(err).To(BeNil())

	//All elements are proper handlers
	length := len(expectedArray)
	for i := 0; i < length; i++ {
		compareHandlers(actual.GetIndex(i), expected.GetIndex(i))
	}

}

// Compares `handler` objects, a common type for `authenticators`, `authorizer`, and `mutator` configurations
// The object consists of two properties: `hander`:string` and `config`:object
func compareHandlers(actual *json.Json, expected *json.Json) {
	//expected.SetPath(
	Expect(actual).To(Equal(expected))

	expectAllowedKeys(actual, "handler", "config")

	expectString(actual, "handler")
	expectObjectOrNil(actual, "config")
}

func expectBoolean(data *json.Json) {
	_, err := data.Bool()
	Expect(err).To(BeNil())
}

func expectString(data *json.Json, attributeName string) {
	errMsg := ""
	_, err := data.Get(attributeName).String()
	if err != nil {
		errMsg = fmt.Sprintf("Cannot convert %s to string. Details: %v", attributeName, err)
	}
	Expect(errMsg).To(BeEmpty())
}

func expectStringArray(data *json.Json, attributeName string) {
	errMsg := ""
	arr, err := data.Get(attributeName).Array()
	if err != nil {
		errMsg = fmt.Sprintf("Cannot convert %s to slice/array. Details: %v", attributeName, err)
	}
	Expect(errMsg).To(BeEmpty())

	length := len(arr)
	for i := 0; i < length; i++ {
		_, ok := arr[i].(string)
		if !ok {
			errMsg = fmt.Sprintf("Cannot convert element of %s array [%d] to string. Details: %v", attributeName, i, err)
		}
		Expect(errMsg).To(BeEmpty())
	}
}

func expectObjectOrNil(data *json.Json, attributeName string) {
	if data.Get("attributeName").Interface() == nil {
		return
	}

	errMsg := ""
	_, err := data.Map()
	if err != nil {
		errMsg = fmt.Sprintf("Cannot convert %s to an object. Details: %v", attributeName, err)
	}
	Expect(errMsg).To(BeEmpty())
}

func expectAllowedKeys(genericMap *json.Json, allowedKeys ...string) {
	aMap, ok := genericMap.Interface().(map[string]interface{})
	Expect(ok).To(BeTrue())
	actualKeys := getKeysOf(aMap)
	for _, v := range actualKeys {
		Expect(allowedKeys).To(ContainElement(v))
	}
}

func expectOnlyKeys(genericMap *json.Json, keys ...string) {
	aMap, ok := genericMap.Interface().(map[string]interface{})
	Expect(ok).To(BeTrue())
	Expect(getKeysOf(aMap)).To(ConsistOf(keys))
}

func getKeysOf(input map[string]interface{}) []string {
	keys := make([]string, len(input))
	i := 0

	for k := range input {
		keys[i] = k
		i++
	}
	return keys
}

// Converts from dynamic client representation to *json.Json.
// "spec" must be a top-level attribute of dynamicObject.
func wrapSpecAsJson(dynamicObject *unstructured.Unstructured) *json.Json {
	//A little trick since go-simplejson doesn't offer a constructor for arbitrary data
	res := json.New()
	res.Set("data", dynamicObject.UnstructuredContent())
	return res.GetPath("data", "spec")
}
