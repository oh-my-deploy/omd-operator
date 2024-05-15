package internal

import (
	"github.com/oh-my-deploy/omd-operator/internal/driver"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OmdManager struct {
	KubeClient    client.Client
	ProgramClient *driver.ProgramClient
}

func NewOmdManager(kube client.Client, schmea *runtime.Scheme) OmdManager {
	programClient := driver.NewProgramClient(kube, schmea)
	return OmdManager{
		KubeClient:    kube,
		ProgramClient: programClient,
	}
}
