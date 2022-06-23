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

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ory/oathkeeper-maester/internal/validation"

	oathkeeperv1alpha1 "github.com/ory/oathkeeper-maester/api/v1alpha1"
	"github.com/ory/oathkeeper-maester/controllers"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme                         = runtime.NewScheme()
	setupLog                       = ctrl.Log.WithName("setup")
)

func init() {

	apiv1.AddToScheme(scheme)
	oathkeeperv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var rulesConfigmapName string
	var rulesConfigmapNamespace string
	var rulesFileName string
	var rulesFilePath string

	var operator controllers.OperatorMode

	controllerCommand := flag.NewFlagSet("controller", flag.ExitOnError)
	sidecarCommand := flag.NewFlagSet("sidecar", flag.ExitOnError)

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	controllerCommand.StringVar(&rulesConfigmapName, "rulesConfigmapName", "oathkeeper-rules", "Name of the Configmap that stores Oathkeeper rules.")
	controllerCommand.StringVar(&rulesConfigmapNamespace, "rulesConfigmapNamespace", "oathkeeper-maester-system", "Namespace of the Configmap that stores Oathkeeper rules.")
	controllerCommand.StringVar(&rulesFileName, "rulesFileName", "access-rules.json", "Name of the key in ConfigMap containing the rules.json")

	sidecarCommand.StringVar(&rulesFilePath, "rulesFilePath", "/etc/config/access-rules.json", "Path to the file with converted Oathkeeper rules")

	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	sideCarMode, err := selectMode(flag.Args(), controllerCommand, sidecarCommand)
	if err != nil {
		setupLog.Error(err, "problem parsing flags")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		// Defaults to "" which means all namespaces
		Namespace: os.Getenv("NAMESPACE"),
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := validateRulesFileName(rulesFileName); err != nil {
		setupLog.Error(err, "Validation error")
		os.Exit(1)
	}

	validationConfig := initValidationConfig()

	if sideCarMode {
		operator = &controllers.FilesOperator{
			Log:           ctrl.Log.WithName("controllers").WithName("Rule"),
			RulesFilePath: rulesFilePath,
		}
	} else {
		operator = &controllers.ConfigMapOperator{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Rule"),
			DefaultConfigMap: types.NamespacedName{
				Name:      rulesConfigmapName,
				Namespace: rulesConfigmapNamespace,
			},
			RulesFileName: rulesFileName,
		}
	}

	ruleReconciler := &controllers.RuleReconciler{
		Client:           mgr.GetClient(),
		Log:              ctrl.Log.WithName("controllers").WithName("Rule"),
		ValidationConfig: validationConfig,
		OperatorMode:     operator,
	}

	if err := ruleReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Rule")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func parseListOrDefault(list string, defaultArr []string, name string) []string {
	if list == "" {
		setupLog.Info(fmt.Sprintf("using default values for %s", name))
		return defaultArr
	}
	return parseList(list)
}

func parseList(list string) []string {
	return removeEmptyStrings(strings.Split(list, ","))
}

func removeEmptyStrings(list []string) []string {
	result := make([]string, 0)
	for _, s := range list {
		ts := strings.TrimSpace(s)
		if ts != "" {
			result = append(result, ts)
		}
	}

	return result
}

func initValidationConfig() validation.Config {
	authenticatorsAvailable := os.Getenv(oathkeeperv1alpha1.AuthenticatorsAvailableEnv)
	authorizersAvailable := os.Getenv(oathkeeperv1alpha1.AuthorizersAvailableEnv)
	mutatorsAvailable := os.Getenv(oathkeeperv1alpha1.MutatorsAvailableEnv)
	return validation.Config{
		AuthenticatorsAvailable: parseListOrDefault(authenticatorsAvailable, oathkeeperv1alpha1.DefaultAuthenticatorsAvailable[:], oathkeeperv1alpha1.AuthenticatorsAvailableEnv),
		AuthorizersAvailable:    parseListOrDefault(authorizersAvailable, oathkeeperv1alpha1.DefaultAuthorizersAvailable[:], oathkeeperv1alpha1.AuthorizersAvailableEnv),
		MutatorsAvailable:       parseListOrDefault(mutatorsAvailable, oathkeeperv1alpha1.DefaultMutatorsAvailable[:], oathkeeperv1alpha1.MutatorsAvailableEnv),
	}
}

func validateRulesFileName(rfn string) error {
	match, _ := regexp.MatchString(oathkeeperv1alpha1.RulesFileNameRegexp, rfn)
	if match {
		return nil
	}
	return fmt.Errorf("rulesFileName: %s is not a valid name", rfn)
}

func selectMode(args []string, controllerCommand *flag.FlagSet, sidecarCommand *flag.FlagSet) (bool, error) {
	sidecar := true
	controller := false

	if len(args) < 1 {
		setupLog.Info("running in controller mode")
		return controller, nil
	}

	switch args[0] {
	case "controller":
		if err := controllerCommand.Parse(args[1:]); err != nil {
			return false, err
		}
		setupLog.Info("running in controller mode")
		return controller, nil
	case "sidecar":
		if err := sidecarCommand.Parse(args[1:]); err != nil {
			return false, err
		}
		setupLog.Info("running in sidecar mode")
		return sidecar, nil
	default:
		return false, fmt.Errorf(`modes "controller" and "sidecar" are supported but got: %s`, args[0])
	}
}
