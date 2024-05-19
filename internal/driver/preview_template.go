package driver

import (
	"context"
	"github.com/oh-my-deploy/omd-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

type PreviewTemplateClient struct {
	KubeClient client.Client
}

func NewPreviewTemplateClient(kube client.Client) *PreviewTemplateClient {
	return &PreviewTemplateClient{
		KubeClient: kube,
	}
}

func (p *PreviewTemplateClient) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	tmpl := &v1alpha1.PreviewTemplate{}
	err := p.KubeClient.Get(ctx, req.NamespacedName, tmpl)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		log.Error(err, "failed to fetch PreviewTemplate")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

//func (p *PreviewTemplateClient) UpsertStatus(ctx context.Context, tmpl *v1alpha1.PreviewTemplate, data string) error {
//	if reflect.DeepEqual(data, tmpl.Spec.Programs) {
//		return nil
//	}
//	tmpl.Status.TemplateData = data
//	return p.KubeClient.Status().Update(ctx, tmpl)
//}
