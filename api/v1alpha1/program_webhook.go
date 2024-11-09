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

package v1alpha1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var programlog = logf.Log.WithName("program-resource")

func (r *Program) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-omd-com-v1alpha1-program,mutating=true,failurePolicy=fail,sideEffects=None,groups=omd.com,resources=programs,verbs=create;update,versions=v1alpha1,name=mprogram.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Program{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Program) Default() {
	programlog.Info("default", "name", r.Name)
	// r.ObjectMeta.Labels = map[string]string{
	// 	"program.kb.io": "true",
	// }
	// r.Annotations = map[string]string{
	// 	"program.kb.io":         "true",
	// 	"program.kb.io/program": r.Name,
	// 	"program.kb.io/owner":   "program-controller",
	// 	"program.kb.io/created": "true",
	// 	"program.kb.io/webhook": "true",
	// }
	r.Annotations["program.kb.io/webhook"] = "true"
	r.Annotations["program.kb.io/created"] = "true"
	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-omd-com-v1alpha1-program,mutating=false,failurePolicy=fail,sideEffects=None,groups=omd.com,resources=programs,verbs=create;update,versions=v1alpha1,name=vprogram.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Program{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Program) ValidateCreate() error {
	programlog.Info("validate create", "name", r.Name)
	if r.Annotations["program.kb.io"] != "true" {
		return errors.New("program.kb.io annotation is not set")
	}
	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Program) ValidateUpdate(old runtime.Object) error {
	programlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Program) ValidateDelete() error {
	programlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
