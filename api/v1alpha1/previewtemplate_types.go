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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PreviewTemplateSpec defines the desired state of PreviewTemplate
type PreviewTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PreviewTemplate. Edit previewtemplate_types.go to remove/update
	//Programs []ProgramSpec `json:"template"`
	Data string `json:"template,omitempty"`
}

// PreviewTemplateStatus defines the observed state of PreviewTemplate
type PreviewTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//TemplateData string `json:"templateData,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PreviewTemplate is the Schema for the previewtemplates API
type PreviewTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PreviewTemplateSpec   `json:"spec,omitempty"`
	Status PreviewTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PreviewTemplateList contains a list of PreviewTemplate
type PreviewTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PreviewTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PreviewTemplate{}, &PreviewTemplateList{})
}
