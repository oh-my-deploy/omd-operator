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
	appv1 "k8s.io/api/apps/v1"
	v2 "k8s.io/api/autoscaling/v2"
	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/apis/networking"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type AppSpec struct {
	Container         v1.Container      `json:"container"`
	Replicas          int32             `json:"replicas,omitempty"`
	AppType           string            `json:"appType,omitempty"` // back, front-spa, front-srr
	PodAnnotations    map[string]string `json:"podAnnotations,omitempty"`
	DeployAnnotations map[string]string `json:"deployAnnotations,omitempty"`
}

type DeploySpec struct {
	Branch string `json:"branch"`
	Path   string `json:"path"`
	Repo   string `json:"repo"`
	Server string `json:"server"`
}

type PodDisruptionBudgetSpec struct {
	// +optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`
	// +optional
	MinAvailable int32 `json:"minAvailable,omitempty"`
	// +optional
	MaxUnavailable int32 `json:"maxUnavailable,omitempty"`
}

type ServiceSpec struct {
	// +kubebuilder:default=false
	Enabled     bool              `json:"enabled"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type IngressSpec struct {
	// +kubebuilder:default=false
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
	Port int32 `json:"port,omitempty"`
}

type ServiceAccountSpec struct {
	// +kubebuilder:default=false
	Create bool `json:"create,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	AutomountServiceAccountToken *bool  `json:"automountServiceAccountToken,omitempty"`
	ServiceAccountName           string `json:"serviceAccountName,omitempty"`
}

type HorizontalPodAutoScalerSpec struct {
	Enabled     bool              `json:"enabled,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	MinReplicas int32             `json:"minReplicas,omitempty"`
	MaxReplicas int32             `json:"maxReplicas,omitempty"`
	Metrics     []v2.MetricSpec   `json:"metrics,omitempty"`
}

type SchedulerSpec struct {
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// +optional
	PodDisruptionBudget *PodDisruptionBudgetSpec `json:"pdb,omitempty"`
	// +optional
	HorizontalPodAutoScaler *HorizontalPodAutoScalerSpec `json:"hpa,omitempty"`
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`
}

// ProgramSpec defines the desired state of Program
type ProgramSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	App            AppSpec            `json:"app,omitempty"`
	Service        ServiceSpec        `json:"service,omitempty"`
	ServiceAccount ServiceAccountSpec `json:"serviceAccount,omitempty"`
	Ingress        IngressSpec        `json:"ingress,omitempty"`
	Scheduler      SchedulerSpec      `json:"scheduler,omitempty"`
	Deploy         DeploySpec         `json:"deploy,omitempty"`
}

// ProgramStatus defines the observed state of Program
type ProgramStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Foo string `json:"foo,omitempty"`
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
	service := v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: p.Name,
			Labels: map[string]string{
				"app": p.Name,
			},
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": p.ObjectMeta.Name,
			},
			Ports: []v1.ServicePort{
				{
					Name: "port",
					Port: p.Spec.App.Container.Ports[0].ContainerPort,
					TargetPort: intstr.IntOrString{
						IntVal: p.Spec.App.Container.Ports[0].ContainerPort,
					},
				},
			},
		},
	}
	if p.Spec.Service.Annotations != nil {
		service.Annotations = p.Spec.Service.Annotations
	}
	return service
}

func (p *Program) ConvertToServiceAccount() v1.ServiceAccount {
	sa := v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: p.Name,
			Labels: map[string]string{
				"app": p.Name,
			},
		},
		AutomountServiceAccountToken: p.Spec.ServiceAccount.AutomountServiceAccountToken,
	}
	if p.Spec.ServiceAccount.Annotations != nil {
		sa.Annotations = p.Spec.ServiceAccount.Annotations
	}
	return sa
}

func (p *Program) ConvertToDeployment() appv1.Deployment {
	deployment := appv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: p.Name,
			Labels: map[string]string{
				"app": p.Name,
			},
		},
		Spec: appv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": p.Name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: p.Name,
					Labels: map[string]string{
						"app": p.Name,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						p.Spec.App.Container,
					},
				},
			},
		},
	}

	if p.Spec.App.DeployAnnotations != nil {
		deployment.Annotations = p.Spec.App.DeployAnnotations
	}
	if p.Spec.App.PodAnnotations != nil {
		deployment.Spec.Template.ObjectMeta.Annotations = p.Spec.App.PodAnnotations
	}
	if p.Spec.App.Replicas != 0 {
		deployment.Spec.Replicas = &p.Spec.App.Replicas
	}
	if p.Spec.Scheduler.NodeSelector != nil {
		deployment.Spec.Template.Spec.NodeSelector = p.Spec.Scheduler.NodeSelector
	}
	if p.Spec.Scheduler.Affinity != nil {
		deployment.Spec.Template.Spec.Affinity = p.Spec.Scheduler.Affinity
	}

	if p.Spec.ServiceAccount.Create {
		deployment.Spec.Template.Spec.ServiceAccountName = p.Name
	} else {
		deployment.Spec.Template.Spec.ServiceAccountName = p.Spec.ServiceAccount.ServiceAccountName
	}

	return deployment
}

func (p *Program) ConvertToIngress() networking.Ingress {
	ingress := networking.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: p.Name,
			Labels: map[string]string{
				"app": p.Name,
			},
		},
		Spec: networking.IngressSpec{
			Rules: []networking.IngressRule{
				{
					Host: p.Spec.Ingress.Rules.Host,
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: p.createIngressPaths(p.Spec.Ingress.Rules.Paths),
						},
					},
				},
			},
		},
	}

	if p.Spec.Ingress.Annotations != nil {
		ingress.Annotations = p.Spec.Ingress.Annotations
	}
	return ingress
}

func (p *Program) createIngressPaths(rules []IngressPath) []networking.HTTPIngressPath {
	networkPaths := make([]networking.HTTPIngressPath, 0)
	for _, rule := range rules {
		newPath := networking.HTTPIngressPath{
			Path: rule.Path,
			Backend: networking.IngressBackend{
				Service: &networking.IngressServiceBackend{
					Name: rule.ServiceName,
					Port: networking.ServiceBackendPort{
						Number: rule.Port,
					},
				},
			},
		}
		networkPaths = append(networkPaths, newPath)
	}
	return networkPaths
}

func (p *Program) ConvertToPdb() policyv1.PodDisruptionBudget {
	pdb := policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: "policy/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        p.Name,
			Annotations: p.Annotations,
			Labels: map[string]string{
				"app": p.Name,
			},
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": p.Name,
				},
			},
		},
	}
	if p.Spec.Scheduler.PodDisruptionBudget.MaxUnavailable != 0 {
		pdb.Spec.MaxUnavailable = &intstr.IntOrString{
			IntVal: p.Spec.Scheduler.PodDisruptionBudget.MaxUnavailable,
		}
	} else {
		pdb.Spec.MinAvailable = &intstr.IntOrString{
			IntVal: p.Spec.Scheduler.PodDisruptionBudget.MinAvailable,
		}
	}
	return pdb
}

func (p *Program) ConvertToHPA() v2.HorizontalPodAutoscaler {
	hpa := v2.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        p.Name,
			Annotations: p.Annotations,
			Labels: map[string]string{
				"app": p.Name,
			},
		},
		Spec: v2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       p.Name,
			},
			MinReplicas: &p.Spec.Scheduler.HorizontalPodAutoScaler.MinReplicas,
			MaxReplicas: p.Spec.Scheduler.HorizontalPodAutoScaler.MaxReplicas,
			Metrics:     p.Spec.Scheduler.HorizontalPodAutoScaler.Metrics,
		},
	}
	return hpa
}
