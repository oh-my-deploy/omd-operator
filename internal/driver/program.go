package driver

import (
	"context"
	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	omdcomv1alpha1 "github.com/oh-my-deploy/omd-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const applicationFinalizer = "argoproj.io/finalizer"

type ProgramClient struct {
	KubeClient client.Client
	Scheme     *runtime.Scheme
}

func NewProgramClient(kubeClient client.Client, schema *runtime.Scheme) *ProgramClient {
	return &ProgramClient{
		KubeClient: kubeClient,
		Scheme:     schema,
	}
}

func (p *ProgramClient) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog := log.FromContext(ctx)
	program := &omdcomv1alpha1.Program{}
	err := p.KubeClient.Get(ctx, req.NamespacedName, program)
	if err != nil {
		if errors.IsNotFound(err) {
			err = p.Delete(ctx, req)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		klog.Error(err, "Failed to fetch Program")
		return ctrl.Result{}, err
	}

	if err = p.SyncArgo(ctx, req, program); err != nil {
		klog.Error(err, err.Error())
		return ctrl.Result{}, err
	}

	if err = p.UpsertProgramStatus(ctx, program); err != nil {
		klog.Error(err, err.Error())
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (p *ProgramClient) Delete(ctx context.Context, req ctrl.Request) error {
	return p.KubeClient.Delete(ctx, &argocdv1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      req.Name,
			Namespace: "argocd",
		},
	})
}

func (p *ProgramClient) UpsertProgramStatus(ctx context.Context, program *omdcomv1alpha1.Program) error {
	if program.Status.Foo == "bar" {
		return nil
	}
	program.Status.Foo = "bar"
	return p.KubeClient.Status().Update(ctx, program)
}

func (p *ProgramClient) SyncArgo(ctx context.Context, req ctrl.Request, program *omdcomv1alpha1.Program) error {
	klog := log.FromContext(ctx)
	klog.Info("Syncing Argo")
	newApp := p.ConvertToArgoCDApplication(program)
	currentApp := &argocdv1alpha1.Application{}
	argoNamespaceName := types.NamespacedName{Namespace: "argocd", Name: program.Name}
	err := p.KubeClient.Get(ctx, argoNamespaceName, currentApp)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("start creating ArgoCD application")
		err = p.KubeClient.Create(ctx, newApp)
		if err != nil {
			klog.Error(err, err.Error())
			return err
		}
		klog.Info("end Creating ArgoCD application")
		return nil
	} else if err != nil {
		klog.Error(err, err.Error())
		return err
	}

	if !reflect.DeepEqual(currentApp.Spec, newApp.Spec) {
		klog.Info("Update start argo")
		err = p.KubeClient.Update(ctx, currentApp)
		if err != nil {
			klog.Error(err, "Failed to update Argo")
			return err
		}
		klog.Info("Updated ArgoCD application.")
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
