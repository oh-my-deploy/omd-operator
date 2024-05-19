package driver

import (
	"context"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ArgoCDClient struct {
	Kubernetes        client.Client
	applicationClient application.ApplicationServiceClient
	projectClient     project.ProjectServiceClient
	clusterClient     cluster.ClusterServiceClient
}

type ArgoDriverConfig struct {
	Address string
	Token   string
}

func NewArgoClient(kube client.Client) *ArgoCDClient {
	return &ArgoCDClient{
		Kubernetes: kube,
	}
}

func (a *ArgoCDClient) DeleteAppSet(ctx context.Context, name string) error {
	currentApp := &argocdv1alpha1.ApplicationSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "argocd",
		},
	}
	return a.Kubernetes.Delete(ctx, currentApp)
}

func (a *ArgoCDClient) DeleteApp(ctx context.Context, name string) error {
	currentApp := &argocdv1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "argocd",
		},
	}
	return a.Kubernetes.Delete(ctx, currentApp)
}

func (a *ArgoCDClient) CreateApp(ctx context.Context, app *argocdv1alpha1.Application) error {
	return a.Kubernetes.Create(ctx, app)
}

func (a *ArgoCDClient) CreateAppSet(ctx context.Context, app *argocdv1alpha1.ApplicationSet) error {
	return a.Kubernetes.Create(ctx, app)
}
