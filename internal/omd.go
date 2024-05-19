package internal

import (
	"github.com/oh-my-deploy/omd-operator/internal/driver"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OmdManager struct {
	KubeClient            client.Client
	ArgoCDClient          *driver.ArgoCDClient
	ProgramClient         *driver.ProgramClient
	PreviewTemplateClient *driver.PreviewTemplateClient
	PreviewClient         *driver.PreviewClient
}

func NewOmdManager(kube client.Client, scheme *runtime.Scheme) OmdManager {
	programClient := driver.NewProgramClient(kube, scheme)
	previewTemplateClient := driver.NewPreviewTemplateClient(kube)
	argoCDClient := driver.NewArgoClient(kube)
	githubClient := driver.NewGithubClient()
	previewClient := driver.NewPreviewClient(kube, githubClient)
	return OmdManager{
		KubeClient:            kube,
		ProgramClient:         programClient,
		PreviewTemplateClient: previewTemplateClient,
		PreviewClient:         previewClient,
		ArgoCDClient:          argoCDClient,
	}
}
