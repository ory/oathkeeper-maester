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
	"os"

	"github.com/go-logr/logr"
	oathkeeperv1alpha1 "github.com/ory/oathkeeper-k8s-controller/api/v1alpha1"
	"github.com/ory/oathkeeper-k8s-controller/internal/validation"

	"github.com/avast/retry-go"
	apiv1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RuleReconciler reconciles a Rule object
type RuleReconciler struct {
	client.Client
	Log              logr.Logger
	RuleConfigmap    types.NamespacedName
	ValidationConfig validation.Config
}

// +kubebuilder:rbac:groups=oathkeeper.ory.sh,resources=rules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=oathkeeper.ory.sh,resources=rules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *RuleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	ctx := context.Background()
	_ = r.Log.WithValues("rule", req.NamespacedName)

	var rule oathkeeperv1alpha1.Rule
	skipValidation := false

	if err := r.Get(ctx, req.NamespacedName, &rule); err != nil {
		if apierrs.IsNotFound(err) {
			skipValidation = true
		} else {
			return ctrl.Result{}, err
		}
	}

	if !skipValidation {
		if err := rule.ValidateWith(r.ValidationConfig); err != nil {
			rule.Status.Validation = &oathkeeperv1alpha1.Validation{}
			rule.Status.Validation.Valid = boolPtr(false)
			e := err.Error()
			rule.Status.Validation.Error = &e
			r.Log.Info("validation error in Rule %s/%s: \"%s\"", rule.Namespace, rule.Name, e)
			if err := r.Update(ctx, &rule); err != nil {
				r.Log.Error(err, "unable to update Rule status")
				return ctrl.Result{}, err
			}
			// continue, as validation can't be fixed by requeuing request and we still have to update the configmap
		} else {
			// rule valid - set the status
			rule.Status.Validation = &oathkeeperv1alpha1.Validation{}
			rule.Status.Validation.Valid = boolPtr(true)
			if err := r.Update(ctx, &rule); err != nil {
				r.Log.Error(err, "unable to update Rule status")
				return ctrl.Result{}, err
			}
		}
	}

	var rulesList oathkeeperv1alpha1.RuleList

	if err := r.List(ctx, &rulesList); err != nil {
		return ctrl.Result{}, err
	}

	oathkeeperRulesJSON, err := rulesList.FilterNotValid().ToOathkeeperRules()
	if err != nil {
		return ctrl.Result{}, err
	}

	if err = r.updateRulesConfigmap(ctx, string(oathkeeperRulesJSON)); err != nil {
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

	if err := r.Get(ctx, r.RuleConfigmap, &oathkeeperRulesConfigmap); err != nil {
		if apierrs.IsNotFound(err) {

			oathkeeperRulesConfigmap = apiv1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      r.RuleConfigmap.Name,
					Namespace: r.RuleConfigmap.Namespace,
				},
				Data: map[string]string{"rules": data},
			}

			err := retry.Do(func() error {
				return r.Create(ctx, &oathkeeperRulesConfigmap)
			},
				retry.Attempts(5),
				retry.DelayType(retry.BackOffDelay))

			if err != nil {
				r.Log.Error(err, "unable to create configmap")
				os.Exit(1)
			}

			return nil
		}

		return err
	}

	oathkeeperRulesConfigmapCopy := oathkeeperRulesConfigmap.DeepCopy()
	oathkeeperRulesConfigmapCopy.Data = map[string]string{"rules": data}

	if err := r.Update(ctx, oathkeeperRulesConfigmapCopy); err != nil {
		return err
	}

	return nil
}

func boolPtr(b bool) *bool {
	return &b
}
