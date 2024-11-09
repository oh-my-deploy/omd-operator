package driver

import (
	"context"
	"reflect"
	"time"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/oh-my-deploy/omd-operator/api/v1alpha1"
	omdcomv1alpha1 "github.com/oh-my-deploy/omd-operator/api/v1alpha1"
	"github.com/oh-my-deploy/omd-operator/internal/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

const applicationFinalizer = "resources-finalizer.argocd.argoproj.io"
const PROGRAM_FINALIZER = "program.omd.com/finalizer"

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
	if err := p.KubeClient.Get(ctx, req.NamespacedName, program); err != nil {
		if kerrors.IsNotFound(err) {
			klog.Info("successful deleted program")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		klog.Error(err, "failed to get program")
		return ctrl.Result{}, err
	}

	if program.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := p.ensureFinalizer(ctx, program); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		return p.handleDeletion(ctx, program)
	}

	if err = p.SyncArgo(ctx, req, program); err != nil {
		if errors.IsConflict(err) {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 15}, nil
		}
		klog.Error(err, err.Error())
		return ctrl.Result{}, err
	}

	if err = p.UpsertProgramStatus(ctx, program); err != nil {
		klog.Error(err, err.Error())
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (p *ProgramClient) Delete(ctx context.Context, program *omdcomv1alpha1.Program) error {
	klog := log.FromContext(ctx)
	klog.Info("Deleting ArgoCD application")
	return p.KubeClient.Delete(ctx, &argocdv1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      program.Name,
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
			// Finalizers are used to ensure that the resources are not deleted by the garbage collector until the finalizer is removed.
			// https://argo-cd.readthedocs.io/en/stable/user-guide/app_deletion/#deletion-using-kubectl
			Finalizers: []string{applicationFinalizer},
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

func (p *ProgramClient) ensureFinalizer(ctx context.Context, program *v1alpha1.Program) error {
	log := ctrllog.FromContext(ctx)
	if !utils.ContainsString(program.GetFinalizers(), PROGRAM_FINALIZER) {
		log.Info("set finalizer in program")
		(*program).SetFinalizers(append(program.GetFinalizers(), PROGRAM_FINALIZER))
		if err := p.KubeClient.Update(ctx, program); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProgramClient) handleDeletion(ctx context.Context, program *v1alpha1.Program) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	if utils.ContainsString(program.GetFinalizers(), PROGRAM_FINALIZER) {
		log.Info("Processing finalize")
		if err := p.Delete(ctx, program); err != nil {
			return ctrl.Result{}, err
		}
		(*program).SetFinalizers(utils.RemoveString(program.GetFinalizers(), PROGRAM_FINALIZER))
		if err := p.KubeClient.Update(ctx, program); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}
