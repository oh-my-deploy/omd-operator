package internal

import (
	"github.com/oh-my-deploy/omd-operator/internal/driver"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OmdManager struct {
	KubeClient    client.Client
	ProgramClient *driver.ProgramClient
}

func NewOmdManager(kube client.Client) OmdManager {
	programClient := driver.NewProgramClient(kube)
	return OmdManager{
		KubeClient:    kube,
		ProgramClient: programClient,
	}
}
