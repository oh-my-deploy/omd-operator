package driver

import (
	"context"
	"errors"
	"github.com/google/go-github/v62/github"
	"net/http"
	"os"
)

type GithubClient struct {
	client *github.Client
}

func NewGithubClient() *GithubClient {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	return &GithubClient{
		client: client,
	}
}

func (g *GithubClient) DispatchWorkflow(ctx context.Context, owner string, repoName string, ref string, inputs map[string]interface{}) error {
	workflowDispatch := github.CreateWorkflowDispatchEventRequest{
		Ref:    ref,
		Inputs: inputs,
	}
	_, err := g.client.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repoName, "preview.yaml", workflowDispatch)
	if err != nil {
		return err
	}
	return nil
}

func (g *GithubClient) CreateOperatorFile(ctx context.Context, owner string, repoName string, data string, path string, previewName string) error {

	opts := &github.RepositoryContentFileOptions{
		Content: []byte(data),
		Branch:  github.String("dev"),
		Message: github.String("create program resource file for  " + previewName),
	}
	sha, err := g.GetOperatorFileSHA(ctx, owner, repoName, path)
	if err != nil {
		return err
	}
	if sha != nil {
		opts.SHA = sha
	}
	_, _, err = g.client.Repositories.CreateFile(ctx, owner, repoName, path, opts)
	return err
}

func (g *GithubClient) GetOperatorFileSHA(ctx context.Context, owner string, repoName string, path string) (*string, error) {
	file, _, res, err := g.client.Repositories.GetContents(ctx, owner, repoName, path, &github.RepositoryContentGetOptions{
		Ref: "dev",
	})
	if res != nil {
		if res.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	if err != nil {
		return nil, err
	}
	return file.SHA, nil
}

func (g *GithubClient) DeleteOperatorFile(ctx context.Context, owner string, repoName string, path string) error {
	opts := &github.RepositoryContentFileOptions{
		Branch:  github.String("dev"),
		Message: github.String("delete program resource file for  " + path),
	}
	sha, err := g.GetOperatorFileSHA(ctx, owner, repoName, path)
	if err != nil {
		return err
	}
	if sha != nil {
		opts.SHA = sha
	}
	_, res, err := g.client.Repositories.DeleteFile(ctx, "oh-my-deploy", "omd-operator-example", path, opts)
	if res != nil {
		if res.StatusCode == http.StatusNotFound {
			return errors.New("resource file not found")
		}
	}
	return err
}

func (g *GithubClient) CreateWorkflowDispatch(ctx context.Context, owner string, repoName string, ref string, inputs map[string]interface{}) error {
	workflowDispatch := github.CreateWorkflowDispatchEventRequest{
		Ref:    "dev",
		Inputs: inputs,
	}
	res, err := g.client.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repoName, "preview.yaml", workflowDispatch)
	if res != nil {
		if res.StatusCode == http.StatusNotFound {
			return errors.New("resource file not found")
		}
	}
	return err
}

func (g *GithubClient) FetchWorkflow(ctx context.Context, owner string, repoName string, workflowID int64) (string, error) {
	workflow, res, err := g.client.Actions.GetWorkflowByID(ctx, owner, repoName, workflowID)
	if res != nil {
		if res.StatusCode == http.StatusNotFound {
			return "", errors.New("resource file not found")
		}
	}
	if err != nil {
		return "", err
	}
	return *workflow.State, nil
}
