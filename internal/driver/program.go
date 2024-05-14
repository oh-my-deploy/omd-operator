package driver

import (
	"context"
	"reflect"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	omdcomv1alpha1 "github.com/oh-my-deploy/omd-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ProgramClient struct {
	KubeClient client.Client
}

func NewProgramClient(kubeClient client.Client) *ProgramClient {
	return &ProgramClient{
		KubeClient: kubeClient,
	}
}

func (p *ProgramClient) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog := log.FromContext(ctx)
	klog.Info("Reconciling start Program")
	program := &omdcomv1alpha1.Program{}
	err := p.KubeClient.Get(ctx, req.NamespacedName, program)
	if err != nil {
		klog.Error(err, err.Error())
		if errors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		return ctrl.Result{}, err
	}
	err = p.SyncArgo(ctx, req, program)
	if err != nil {
		klog.Error(err, err.Error())
		return ctrl.Result{}, err
	}
	klog.Info("Reconciling end Program")
	return ctrl.Result{}, nil
}

func (p *ProgramClient) SyncArgo(ctx context.Context, req ctrl.Request, program *omdcomv1alpha1.Program) error {
	klog := log.FromContext(ctx)
	klog.Info("Syncing Argo")
	newApp := p.ConvertToArgoCDApplication(program)
	app := &argocdv1alpha1.Application{}
	err := p.KubeClient.Get(ctx, req.NamespacedName, app)
	if app != nil && !reflect.DeepEqual(app.Spec, newApp.Spec) {
		err = p.KubeClient.Update(ctx, newApp)
		if err != nil {
			klog.Info("Failed to update ArgoCD application.")
			return err
		}
		klog.Info("Updated ArgoCD application.")
	} else if app == nil {
		klog.Info("start creating ArgoCD application")
		err = p.KubeClient.Create(ctx, newApp)
		if err != nil {
			klog.Error(err, err.Error())
			return err
		}
		klog.Info("end creating ArgoCD application")
	} else {
		klog.Info("Synced ArgoCD application")
	}
	return nil
}

func (p *ProgramClient) ConvertToArgoCDApplication(program *omdcomv1alpha1.Program) *argocdv1alpha1.Application {
	return &argocdv1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      program.Name,
			Namespace: "argocd",
		},
		Spec: argocdv1alpha1.ApplicationSpec{
			Destination: argocdv1alpha1.ApplicationDestination{
				Server:    program.Spec.Deploy.Server,
				Namespace: "default",
			},
			Project: "default",
			SyncPolicy: &argocdv1alpha1.SyncPolicy{
				Automated: &argocdv1alpha1.SyncPolicyAutomated{
					Prune:    true,
					SelfHeal: true,
				},
			},
			Source: &argocdv1alpha1.ApplicationSource{
				RepoURL:        program.Spec.Deploy.Repo,
				Path:           program.Spec.Deploy.Path,
				TargetRevision: program.Spec.Deploy.Branch,
			},
		},
	}
}
