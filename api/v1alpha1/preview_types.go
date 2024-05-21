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
	"encoding/json"
	"errors"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type GitStatus string

const (
	GIT_STATUS_PENDING    GitStatus = "PENDING"
	GIT_STATUS_PROCESSING GitStatus = "PROCESSING"
	GIT_STATUS_FAILED     GitStatus = "FAILED"
	GIT_STATUS_SUCCESS    GitStatus = "SUCCESS"
)

type ActionArg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PreviewProgramSpec struct {
	Program        ProgramSpec `json:"spec"`
	RepositoryName string      `json:"repositoryName"`
	Branch         string      `json:"branch"`
	ActionArgs     []ActionArg `json:"actionArgs,omitempty"`
}

// PreviewSpec defines the desired state of Preview
type PreviewSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Programs        []PreviewProgramSpec `json:"programs,omitempty"`
	PreviewTemplate PreviewTemplateRef   `json:"template,omitempty"`
}

type PreviewTemplateRef struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params,omitempty"`
}

//type Params struct {
//	Key   string `json:"key,omitempty" structs:"key,omitempty"`
//	Value string `json:"value,omitempty" structs:"value,omitempty"`
//}

// PreviewStatus defines the observed state of Preview
type PreviewStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	TemplateSpec            string                    `json:"templateSpec,omitempty"`
	DeployedStatus          []DeployedStatus          `json:"deployedStatus,omitempty"`
	CreatePreviewSpecStatus []CreatePreviewSpecStatus `json:"createPreviewSpecStatus,omitempty"`
}
type DeployedStatus struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	ActionID string `json:"actionID"`
}

type CreatePreviewSpecStatus struct {
	Message   string    `json:"message"`
	GitStatus GitStatus `json:"gitStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Preview is the Schema for the previews API
type Preview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PreviewSpec   `json:"spec,omitempty"`
	Status PreviewStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PreviewList contains a list of Preview
type PreviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Preview `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Preview{}, &PreviewList{})
}

func (p *Preview) ConvertDeployData(paths []string) ([]apiextensionsv1.JSON, error) {
	jsonArr := []apiextensionsv1.JSON{}
	var deployedStatuses []DeployedStatus
	for idx, program := range p.Spec.Programs {
		b, err := json.Marshal(map[string]interface{}{
			"repository-name": program.RepositoryName,
			"target-path":     paths[idx],
		})
		deployedStatuses = append(deployedStatuses, DeployedStatus{
			Name: program.RepositoryName,
			Path: paths[idx],
		})
		if err != nil {
			return nil, errors.Join(err, errors.New("failed to marshal deploy"))
		}
		p.Spec.Programs[idx].Program.Deploy.Path = paths[idx]
		jsonArr = append(jsonArr, apiextensionsv1.JSON{Raw: b})
	}
	p.Status.DeployedStatus = deployedStatuses
	return jsonArr, nil
}
