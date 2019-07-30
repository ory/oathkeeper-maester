package integration

import (
	"errors"
	"fmt"
	"os"
	"time"

	json "github.com/bitly/go-simplejson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	oathkeeperv1alpha1 "github.com/ory/oathkeeper-maester/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespaceName                 = "test-namespace" //TODO: Randomize?
	defaultTargetMapNamespace     = "oathkeeper-maester-system"
	defaultTargetMapName          = "oathkeeper-rules"
	maxRetriesWaitingForConfigMap = 35 //Twice as max registered on my machine.
	maxRetriesWaitingForRule      = 10
	rulesFileName                 = "access-rules.json"
)

var (
	k8sClient        *kubernetes.Clientset       = getK8sClientOrDie()
	k8sDynamicClient dynamic.Interface           = getK8sDynamicClientOrDie()
	ruleResource     schema.GroupVersionResource = schema.GroupVersionResource{Group: oathkeeperv1alpha1.GroupVersion.Group, Version: oathkeeperv1alpha1.GroupVersion.Version, Resource: "rules"}
)

var _ = BeforeSuite(func() {
	//Create namespace
	_, err := createNamespace(namespaceName, k8sClient)
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	//Delete namespace
	err := deleteNamespace(namespaceName, k8sClient)
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("Oathkeeper controller", func() {

	Context("should manage rules as ConfigMap entries", func() {

		It("in happy path scenario", func() {

			By("create a valid ConfigMap from a single Rule")
			//Given
			rule, err := getRule(getRule1Json())
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())

			//When
			_, createErr := ensureRule(rule)
			Expect(createErr).To(BeNil())

			//Then
			rulesArray, validateErr := validateConfigMapContains(rule)
			Expect(validateErr).To(BeNil())
			expectRuleCount(rulesArray, 1)

			By("add an entry to the ConfigMap after adding another Rule")
			//Given
			rule, err = getRule(getRule2Json())
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())

			//When
			_, createErr = ensureRule(rule)
			Expect(createErr).To(BeNil())

			//Then
			rulesArray, validateErr = validateConfigMapContains(rule)
			Expect(validateErr).To(BeNil())
			expectRuleCount(rulesArray, 2)

			By("update a ConfigMap entry after Rule update")
			//Given
			rule, getErr := k8sDynamicClient.Resource(ruleResource).Namespace(namespaceName).Get("test-rule-2", metav1.GetOptions{})
			Expect(getErr).To(BeNil())

			//When
			updateRule(rule)
			_, updateErr := k8sDynamicClient.Resource(ruleResource).Namespace(namespaceName).Update(rule, metav1.UpdateOptions{})
			Expect(updateErr).To(BeNil())

			//Allow for some processing time
			time.Sleep(5 * time.Second)

			//Then
			rulesArray, validateErr = validateConfigMapContains(rule)
			Expect(validateErr).To(BeNil())
			expectRuleCount(rulesArray, 2)
			By("delete a ConfigMap entry after Rule delete")
			//When
			deleteErr := k8sDynamicClient.Resource(ruleResource).Namespace(namespaceName).Delete("test-rule-2", &metav1.DeleteOptions{})
			Expect(deleteErr).To(BeNil())

			//Allow for some processing time
			time.Sleep(3 * time.Second)

			//Then
			rule, err = getRule(getRule1Json())
			Expect(err).To(BeNil())

			rulesArray, validateErr = validateConfigMapContains(rule)
			Expect(validateErr).To(BeNil())
			expectRuleCount(rulesArray, 1)

			By("delete last ConfigMap entry after all Rules are deleted")
			//When
			deleteErr = k8sDynamicClient.Resource(ruleResource).Namespace(namespaceName).Delete("test-rule-1", &metav1.DeleteOptions{})
			Expect(deleteErr).To(BeNil())

			//Allow for some processing time
			time.Sleep(3 * time.Second)

			//Then
			emptyMap, err := getTargetMap()
			Expect(err).To(BeNil())
			Expect(emptyMap.Data[rulesFileName]).To(Equal("[]"))
		})
	})
})

//Converts Rule CRD instance to a *json.Json representation of an entry that should be created in the ConfigMap.
//At the moment CRD structure and ConfigMap entry structure is very similar, we can use that to avoid separate json file with "expected" data.
func toExpected(rule *unstructured.Unstructured) *json.Json {
	expected := wrapSpecAsJson(rule)

	expectedRuleId := rule.GetName() + "." + namespaceName
	expected.Set("id", expectedRuleId)

	upstream := expected.GetPath("upstream")

	//preserveHost in Rule => preserve_host in ConfigMap
	upstream.Set("preserve_host", upstream.Get("preserveHost").Interface())
	upstream.Del("preserveHost")

	//The controller always seserializes preserve_host, nil in Rule gets translated to false.
	if upstream.Get("preserve_host").Interface() == nil {
		upstream.Set("preserve_host", false)
	}

	return expected
}

func updateRule(rule *unstructured.Unstructured) {
	emptyMap := map[string]interface{}{}
	ruleJson := wrapSpecAsJson(rule)
	ruleJson.SetPath([]string{"upstream", "preserveHost"}, true)
	ruleJson.SetPath([]string{"upstream", "url"}, "https://xyz.org")
	ruleJson.SetPath([]string{"match", "url"}, "https://fedcba.com")
	ruleJson.SetPath([]string{"authorizer", "handler"}, "deny")
	ruleJson.SetPath([]string{"authorizer", "config"}, emptyMap)
}

func expectRuleCount(rulesArray *json.Json, cnt int) {
	Expect(rulesArray.Interface()).To(HaveLen(cnt))
}

func getK8sDynamicClientOrDie() dynamic.Interface {
	var kubeconfig = os.Getenv("KUBECONFIG")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("\nError building K8s config: %s\n", err.Error())
		os.Exit(1)
	}

	// create the client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("\nError creating K8s dynamic client: %s\n", err.Error())
		os.Exit(1)
	}

	return client
}

func getK8sClientOrDie() *kubernetes.Clientset {
	var kubeconfig = os.Getenv("KUBECONFIG")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("\nError building K8s config: %s\n", err.Error())
		os.Exit(1)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("\nError creating K8s clientset: %s\n", err.Error())
		os.Exit(1)
	}

	return clientset
}

func toNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.NamespaceSpec{},
	}
}

func deleteNamespace(name string, k8sClient *kubernetes.Clientset) error {

	worker := func() (interface{}, error) {
		err := k8sClient.CoreV1().Namespaces().Delete(name, nil)
		return nil, err
	}

	const maxRetries = 2
	_, err := withRetries(maxRetries, time.Second*1, onRetryLogMsg("retry deleting namespace"), worker)
	return err
}

func createNamespace(name string, k8sClient *kubernetes.Clientset) (*corev1.Namespace, error) {

	testNamespace := toNamespace(name)
	return k8sClient.CoreV1().Namespaces().Create(testNamespace)
}

func getRule(json string) (*unstructured.Unstructured, error) {
	res := unstructured.Unstructured{}

	err := res.UnmarshalJSON([]byte(json))

	if err != nil {
		return nil, err
	}
	return &res, nil
}

func getTargetMap() (*v1.ConfigMap, error) {
	return k8sClient.CoreV1().ConfigMaps(getTargetMapNamespace()).Get(getTargetMapName(), metav1.GetOptions{})
}

//Entry point for validation
//Returns parsed rules array as json.Json object
func validateConfigMapContains(sourceRule *unstructured.Unstructured) (*json.Json, error) {

	//It's a copy!
	expectedRule := toExpected(sourceRule.DeepCopy())

	expectedRuleId, ok := expectedRule.Get("id").Interface().(string)
	if !ok {
		return nil, errors.New("Can't find \"id\" in expectedRule")
	}

	workerFunc := func() (interface{}, error) {

		//Fetch target ConfigMap
		targetMap, err := getTargetMap()
		if err != nil {
			return nil, err
		}

		//Parse data from ConfigMap
		jsonString := targetMap.Data[rulesFileName]
		if jsonString == "" || jsonString == "null" {
			return nil, errors.New("No rules in ConfigMap")
		}
		jsonRules, err := json.NewJson([]byte(jsonString))
		if err != nil {
			return nil, err
		}

		//Look for the rule
		actualRule, err := findRule(jsonRules, expectedRuleId)
		if err != nil {
			return nil, err
		}

		//Validate
		validateRuleEquals(actualRule, expectedRule)
		return jsonRules, nil
	}

	//Execute with retries
	res, err := withRetries(maxRetriesWaitingForConfigMap, time.Second*3, onRetryLogMsg("retry validating ConfigMap"), workerFunc)
	if err != nil {
		return nil, err
	}

	rulesArray, ok := res.(*json.Json)
	if !ok {
		return nil, errors.New("Can't convert to rules array")
	}

	return rulesArray, nil
}

//Creates a Rule and ensures it can be read back from the cluster.
func ensureRule(rule *unstructured.Unstructured) (*unstructured.Unstructured, error) {

	//Create
	_, createErr := k8sDynamicClient.Resource(ruleResource).Namespace(namespaceName).Create(rule, metav1.CreateOptions{})
	Expect(createErr).To(BeNil())

	//Wait until Rule can be read from the cluster.
	workerFunc := func() (interface{}, error) {
		return k8sDynamicClient.Resource(ruleResource).Namespace(namespaceName).Get(rule.GetName(), metav1.GetOptions{})
	}

	//Execute with retries
	res, err := withRetries(maxRetriesWaitingForRule, time.Second*3, onRetryLogMsg("retry getting new Rule"), workerFunc)
	if err != nil {
		return nil, err
	}

	createdRule, ok := res.(*unstructured.Unstructured)
	if !ok {
		return nil, errors.New("Can't convert to a Rule")
	}

	return createdRule, nil
}

//Finds and returns a rule with specified id
//actualRules must be an array of rules (sadly go-simplejson doesn't offer such type)
func findRule(actualRules *json.Json, expectedRuleId string) (*json.Json, error) {

	//Can it be done simpler? go-simplejson doesn't seem to offer any generic way of iterating over array.
	testIfArray, err := actualRules.Array()
	if err != nil {
		return nil, err
	}

	ruleCount := len(testIfArray)

	for i := 0; i < ruleCount; i++ {
		rule := actualRules.GetIndex(i)
		ruleId, err := rule.Get("id").String()
		if err != nil {
			return nil, err
		}
		if ruleId == expectedRuleId {
			return rule, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Rule with id %s not found", expectedRuleId))
}

func getTargetMapNamespace() string {
	res := os.Getenv("TARGET_MAP_NAMESPACE")
	if res == "" {
		res = defaultTargetMapNamespace
	}

	return res
}

func getTargetMapName() string {
	res := os.Getenv("TARGET_MAP_NAME")

	if res == "" {
		res = defaultTargetMapName
	}

	return res
}

func getRule1Json() string {
	return `{
		"apiVersion": "oathkeeper.ory.sh/v1alpha1",
		"kind": "Rule",
		"metadata": {
			"name": "test-rule-1"
	    },
		"spec": {
			"match": {
				"methods": ["GET", "POST"],
				"url": "http://gh.ij"
			},
			"upstream": {
			    "preserveHost": false,
			    "url": "http://abc.def"
			},
			"authenticators": [
				{"handler": "anonymous"}
			],
			"authorizer": {
				"handler": "allow"
			},
			"mutator": {
				"handler": "header",
				"config": {
				    "headers": {
						"X-User": "{{ print .Subject }}",
						"X-Some-Arbitrary-Data": "{{ print .Extra.some.arbitrary.data }}"
				    }
				}
			}
		}
	}
	`
}

func getRule2Json() string {
	return `{
		"apiVersion": "oathkeeper.ory.sh/v1alpha1",
		"kind": "Rule",
		"metadata": {
			"name": "test-rule-2"
	    },
		"spec": {
			"match": {
				"methods": ["POST", "PUT"],
				"url": "http://xyz.com"
			},
			"upstream": {
			    "url": "http://abcde.fgh"
			},
			"authenticators": [
				{
					"handler": "oauth2_client_credentials",
					"config": {
						"required_scope": ["scope-a", "scope-b"]
					}
				},
				{"handler": "anonymous"}
			],
			"authorizer": {
			    "handler": "keto_engine_acp_ory",
				"config": {
				    "required_action": "my:action:1234",
					"required_resource": "my:resource:foobar:foo:1234"
				}
			},
			"mutator": {
			    "handler": "id_token",
			    "config": {
				    "aud": ["audience1", "audience2"]
			    }
			}
		}
	}
	`
}
