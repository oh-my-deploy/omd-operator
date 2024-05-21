package driver

import (
	"context"
	"errors"
	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/oh-my-deploy/omd-operator/api/v1alpha1"
	"github.com/oh-my-deploy/omd-operator/internal/utils"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

type PreviewClient struct {
	KubeClient   client.Client
	GithubClient *GithubClient
}

func NewPreviewClient(kubeClient client.Client, GithubClient *GithubClient) *PreviewClient {
	return &PreviewClient{
		KubeClient:   kubeClient,
		GithubClient: GithubClient,
	}
}

func (p *PreviewClient) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	preview := &v1alpha1.Preview{}
	log.Info("start preview reconcile")
	var result ctrl.Result
	if err := p.KubeClient.Get(ctx, req.NamespacedName, preview); err != nil {
		if kerrors.IsNotFound(err) {
			log.Info("preview not found", "name ", req.Name)
			err := p.Delete(ctx, req.Name)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		log.Error(err, "failed to get preview")
		return ctrl.Result{}, err
	}
	if err := p.ParsingPreviewTemplate(ctx, preview); err != nil {
		log.Error(err, "failed to parse preview template")
		return ctrl.Result{}, err
	}
	// TODO: create application set using argocd
	if err := p.UpsertApplicationSet(ctx, preview); err != nil {
		log.Error(err, "failed to upsert application set")
		return ctrl.Result{}, err
	}
	////TODO: create yaml data, upload deploy repo
	if result, err := p.CreatePreviewContent(ctx, preview); err != nil {
		log.Error(err, "failed to upsert deploy repo")
		if result.Requeue {
			return result, nil
		}
		return result, err
	}
	log.Info("finish preview reconcile")
	//TODO: build container image using github action for preview

	//log.Info("end preview reconcile", "result", result)
	return result, nil
}

func (p *PreviewClient) ParsingPreviewTemplate(ctx context.Context, preview *v1alpha1.Preview) error {
	log := ctrllog.FromContext(ctx)
	var isReset bool
	log.Info("start ParsingPreviewTemplate")
	if len((*preview).Spec.Programs) != 0 || preview.Spec.PreviewTemplate.Name == "" {
		return nil
	}
	previewTemplate := &v1alpha1.PreviewTemplate{}
	var previewSpecs []v1alpha1.PreviewProgramSpec
	if err := p.KubeClient.Get(ctx, types.NamespacedName{Name: (*preview).Spec.PreviewTemplate.Name, Namespace: (*preview).Namespace}, previewTemplate); err != nil {
		log.Error(err, "failed to get previewTemplate")
		return err
	}
	parsedData, err := utils.ParsingPreviewTemplate(previewTemplate.Spec.Data, (*preview).Spec.PreviewTemplate.Params)
	if err != nil {
		log.Error(err, "failed to parse previewTemplate")
		return err
	}
	if err := utils.ConvertToObj(parsedData, &previewSpecs); err != nil {
		log.Error(err, "failed to convert to obj")
		return err
	}
	if !reflect.DeepEqual(preview.Status.TemplateSpec, parsedData) {
		(*preview).Status.TemplateSpec = parsedData
		if err := p.KubeClient.Status().Update(ctx, preview); err != nil {
			log.Error(err, "failed to update preview status")
			return err
		}
		isReset = true
	}
	(*preview).Spec.Programs = previewSpecs
	if len((*preview).Status.CreatePreviewSpecStatus) == 0 || isReset == true {
		log.Info("reset status about createPreviewStatus")
		(*preview).Status.CreatePreviewSpecStatus = make([]v1alpha1.CreatePreviewSpecStatus, len(previewSpecs))
	}
	log.Info("end ParsingPreviewTemplate")
	return nil
}

func (p *PreviewClient) Delete(ctx context.Context, name string) error {
	currentApp := &argocdv1alpha1.ApplicationSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "argocd",
		},
	}
	_ = p.KubeClient.Delete(ctx, currentApp)

	_ = p.KubeClient.Delete(ctx, &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
	})
	return nil
}

func (p *PreviewClient) UpsertApplicationSet(ctx context.Context, preview *v1alpha1.Preview) error {
	log := ctrllog.FromContext(ctx)
	log.Info("start upsert application set")
	currentApp := &argocdv1alpha1.ApplicationSet{}
	newApp := p.GenerateApplicationSet(preview)
	if err := p.KubeClient.Get(ctx, types.NamespacedName{Name: preview.Name, Namespace: "argocd"}, currentApp); err != nil {
		if kerrors.IsNotFound(err) {
			paths := utils.RandomStringLists(len(preview.Spec.Programs))
			generatorData, err := (*preview).ConvertDeployData(paths)
			if err != nil {
				log.Error(err, "failed to convert deploy data")
				return err
			}
			newApp.Spec.Generators = []argocdv1alpha1.ApplicationSetGenerator{{
				List: &argocdv1alpha1.ListGenerator{
					Elements: generatorData,
				},
			}}
			if err = p.KubeClient.Create(ctx, &newApp); err != nil {
				log.Error(err, "failed to create application set")
				return err
			}
			return client.IgnoreNotFound(err)
		}
		return err
	}
	if !reflect.DeepEqual(currentApp.Spec.Template, newApp.Spec.Template) {
		log.Info("update application set")
		currentApp.Spec.Template = newApp.Spec.Template
		if err := p.KubeClient.Update(ctx, currentApp); err != nil {
			log.Error(err, "failed to update application set")
			return err
		}
	}
	log.Info("end upsert application set")
	return nil
}

func (p *PreviewClient) GenerateApplicationSet(preview *v1alpha1.Preview) argocdv1alpha1.ApplicationSet {
	return argocdv1alpha1.ApplicationSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      preview.Name,
			Namespace: "argocd",
		},
		Spec: argocdv1alpha1.ApplicationSetSpec{
			SyncPolicy: &argocdv1alpha1.ApplicationSetSyncPolicy{
				PreserveResourcesOnDeletion: true,
			},
			Template: argocdv1alpha1.ApplicationSetTemplate{
				ApplicationSetTemplateMeta: argocdv1alpha1.ApplicationSetTemplateMeta{
					Name:      preview.Name + "-" + "{{repository-name}}",
					Namespace: "argocd",
				},
				Spec: argocdv1alpha1.ApplicationSpec{
					Source: &argocdv1alpha1.ApplicationSource{
						RepoURL:        "git@github.com:oh-my-deploy/omd-operator-example.git",
						Path:           "{{target-path}}",
						TargetRevision: "dev",
					},
					Destination: argocdv1alpha1.ApplicationDestination{
						Namespace: preview.Name,
						Server:    "https://kubernetes.default.svc",
					},
					SyncPolicy: &argocdv1alpha1.SyncPolicy{
						Automated: &argocdv1alpha1.SyncPolicyAutomated{
							Prune:    true,
							SelfHeal: true,
						},
						SyncOptions: argocdv1alpha1.SyncOptions{
							"CreateNamespace=true",
						},
					},
				},
			},
		},
	}
}

func (p *PreviewClient) CreatePreviewContent(ctx context.Context, preview *v1alpha1.Preview) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	var issuedError bool
	currentPreview := preview.DeepCopy()
	for idx, status := range preview.Status.DeployedStatus {
		createFileStatus := preview.Status.CreatePreviewSpecStatus[idx]
		previewProgram := &v1alpha1.Program{
			TypeMeta: v1.TypeMeta{
				APIVersion: "omd.com/v1alpha1",
				Kind:       "Program",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: preview.Name + "-" + status.Name,
			},
			Spec: preview.Spec.Programs[idx].Program,
		}

		switch createFileStatus.GitStatus {
		case v1alpha1.GIT_STATUS_SUCCESS:
			continue
		default:
			log.Info("start create preview content")
			data, err := utils.ConvertToYaml(previewProgram)
			if err != nil {
				log.Error(err, "failed to convert to yaml")
				(*preview).Status.CreatePreviewSpecStatus[idx].Message = err.Error()
				issuedError = true
				continue
			}
			if err := p.GithubClient.CreateOperatorFile(ctx, "oh-my-deploy", "omd-operator-example", data, status.Path+"/operator.yaml", preview.Name); err != nil {
				log.Error(err, "failed to create operator file")
				(*preview).Status.CreatePreviewSpecStatus[idx].Message = err.Error()
				issuedError = true
			}
			if issuedError {
				(*preview).Status.CreatePreviewSpecStatus[idx].GitStatus = v1alpha1.GIT_STATUS_FAILED
			} else {
				(*preview).Status.CreatePreviewSpecStatus[idx].GitStatus = v1alpha1.GIT_STATUS_SUCCESS
			}
			log.Info("end create preview content")
		}
	}
	if !reflect.DeepEqual(currentPreview.Status.CreatePreviewSpecStatus, (*preview).Status.CreatePreviewSpecStatus) {
		if err := p.KubeClient.Status().Update(ctx, preview); err != nil {
			log.Error(err, "failed to update preview status")
			return ctrl.Result{}, err
		}
	}
	if issuedError {
		log.Info("issued error")
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, errors.New("failed to issue preview content")
	}
	return ctrl.Result{}, nil
}
