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
	v2 "k8s.io/api/autoscaling/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type AppSpec struct {
	Image         string            `json:"image"`
	ContainerPort int32             `json:"containerPort"`
	Replicas      *int32            `json:"replicas,omitempty"`
	AppType       string            `json:"appType,omitempty"` // back, front-spa, front-srr
	Annotations   map[string]string `json:"annotations,omitempty"`
	Probe         ProbeSpec         `json:"probe,omitempty"`
}

type DeploySpec struct {
	Branch string `json:"branch"`
	Path   string `json:"path"`
	Repo   string `json:"repo"`
	Server string `json:"server"`
}

type ProbeSpec struct {
	Startup   *v1.Probe `json:"startup,omitempty"`
	Liveness  *v1.Probe `json:"liveness,omitempty"`
	Readiness *v1.Probe `json:"readiness,omitempty"`
}

type PodDisruptionBudgetSpec struct {
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// +optional
	MinAvailable *int32 `json:"minAvailable,omitempty"`
	// +optional
	MaxUnavailable *int32 `json:"maxUnavailable,omitempty"`
}

type ServiceSpec struct {
	Enabled     bool              `json:"enabled"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type IngressSpec struct {
	Enabled     bool              `json:"enabled"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Rules       IngressRulesSpec  `json:"rules"`
}

type IngressRulesSpec struct {
	Host string `json:"host,omitempty"`
	// +listType=atomic
	Paths []IngressPath `json:"paths"`
}

type IngressPath struct {
	// +optional
	Path string `json:"path,omitempty"`
	// +optional
	ServiceName string `json:"serviceName,omitempty"`
	// +optional
	Port *int32 `json:"port,omitempty"`
}

type ServiceAccountSpec struct {
	Create *bool `json:"create,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	AutomountServiceAccountToken *bool  `json:"automountServiceAccountToken,omitempty"`
	ServiceAccountName           string `json:"serviceAccountName,omitempty"`
}

type HorizontalPodAutoScalerSpec struct {
	Enabled     *bool             `json:"enabled,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	MinReplicas *int32            `json:"minReplicas,omitempty"`
	MaxReplicas *int32            `json:"maxReplicas,omitempty"`
	Metrics     []v2.MetricSpec   `json:"metrics,omitempty"`
}

type SchedulerSpec struct {
	NodeSelector            map[string]string            `json:"nodeSelector,omitempty"`
	PodDisruptionBudget     *PodDisruptionBudgetSpec     `json:"pdb,omitempty"`
	HorizontalPodAutoScaler *HorizontalPodAutoScalerSpec `json:"hpa,omitempty"`
	Affinity                *v1.Affinity                 `json:"affinity,omitempty"`
}

// ProgramSpec defines the desired state of Program
type ProgramSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Foo is an example field of Program. Edit program_types.go to remove/update
	Foo            string              `json:"foo,omitempty"`
	App            *AppSpec            `json:"app,omitempty"`
	Service        *ServiceSpec        `json:"service,omitempty"`
	ServiceAccount *ServiceAccountSpec `json:"serviceAccount,omitempty"`
	Ingress        *IngressSpec        `json:"ingress,omitempty"`
	Scheduler      *SchedulerSpec      `json:"scheduler,omitempty"`
	Deploy         *DeploySpec         `json:"deploy,omitempty"`
}

// ProgramStatus defines the observed state of Program
type ProgramStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Program is the Schema for the programs API
type Program struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProgramSpec   `json:"spec,omitempty"`
	Status ProgramStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProgramList contains a list of Program
type ProgramList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Program `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Program{}, &ProgramList{})
}

func (p *Program) ConvertToService() v1.Service {

	return v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: p.Name,
			Labels: map[string]string{
				"app": p.Name,
			},
			Annotations: p.Spec.Service.Annotations,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": p.ObjectMeta.Name,
			},
			Ports: []v1.ServicePort{
				{
					Name: "port",
					Port: p.Spec.App.ContainerPort,
					TargetPort: intstr.IntOrString{
						IntVal: p.Spec.App.ContainerPort,
					},
				},
			},
		},
	}
}

func (p *Program) ConvertToServiceAccount() v1.ServiceAccount {
	return v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        p.Name,
			Annotations: p.Spec.ServiceAccount.Annotations,
			Labels: map[string]string{
				"app": p.Name,
			},
		},
		AutomountServiceAccountToken: p.Spec.ServiceAccount.AutomountServiceAccountToken,
	}
}
