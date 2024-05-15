/*
Copyright 2024.

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
	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"

	omdcomv1alpha1 "github.com/oh-my-deploy/omd-operator/api/v1alpha1"
	"github.com/oh-my-deploy/omd-operator/internal"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ProgramReconciler reconciles a Program object
type ProgramReconciler struct {
	OmdManager internal.OmdManager
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=omd.com,resources=programs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=omd.com,resources=programs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=omd.com,resources=programs/finalizers,verbs=update
//+kubebuilder:rbac:groups=argoproj.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;create;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Program object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ProgramReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.OmdManager.ProgramClient.Reconcile(ctx, req)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProgramReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&omdcomv1alpha1.Program{}).
		Owns(&argocdv1alpha1.Application{}).
		Complete(r)
}
