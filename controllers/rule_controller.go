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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	oathkeeperv1alpha1 "github.com/ory/oathkeeper-k8s-controller/api/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RuleReconciler reconciles a Rule object
type RuleReconciler struct {
	client.Client
	Log           logr.Logger
	RuleConfigmap types.NamespacedName
}

// +kubebuilder:rbac:groups=oathkeeper.ory.sh,resources=rules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=oathkeeper.ory.sh,resources=rules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *RuleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	ctx := context.Background()
	_ = r.Log.WithValues("rule", req.NamespacedName)
	rulesList := &oathkeeperv1alpha1.RuleList{}

	err := r.List(ctx, rulesList)
	if err != nil {
		return ctrl.Result{}, err
	}

	oathkeeperRulesJSON, err := rulesList.ToOathkeeperRules()
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.updateRulesConfigmap(ctx, string(oathkeeperRulesJSON))
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *RuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oathkeeperv1alpha1.Rule{}).
		Owns(&apiv1.ConfigMap{}).
		Complete(r)
}

func (r *RuleReconciler) updateRulesConfigmap(ctx context.Context, data string) error {

	var oathkeeperRulesConfigmap apiv1.ConfigMap

	err := r.Get(ctx, r.RuleConfigmap, &oathkeeperRulesConfigmap)
	if err != nil {
		return err
	}

	oathkeeperRulesConfigmapCopy := oathkeeperRulesConfigmap.DeepCopy()
	oathkeeperRulesConfigmapCopy.Data = map[string]string{"rules": data}

	err = r.Update(ctx, oathkeeperRulesConfigmapCopy)
	if err != nil {
		return err
	}

	return nil
}
