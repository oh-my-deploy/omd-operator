package internal

import "sigs.k8s.io/controller-runtime/pkg/client"

type OmdManager struct {
	KubeClient client.Client
}

func NewOmdManager(kube client.Client) OmdManager {
	return OmdManager{
		KubeClient: kube,
	}
}
