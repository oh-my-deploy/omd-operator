package driver

import (
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
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

func NewArgoClient(cfg ArgoDriverConfig) ArgoCDClient {
	return ArgoCDClient{}
}
