package driver_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/oh-my-deploy/omd-operator/internal/config"
	"github.com/oh-my-deploy/omd-operator/internal/driver"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var path string
var client *driver.GithubClient
var owner string
var repoName string

func TestMain(m *testing.M) {
	path = "test-code/operator.yaml"
	owner = "oh-my-deploy"
	repoName = "omd-operator-example"
	config.LoadEnv("../../")
	client = driver.NewGithubClient()
	code := m.Run()
	os.Exit(code)
}

func TestGithubClient_CreateOperatorFile(t *testing.T) {
	randomData := uuid.New().String()
	t.Run("CreateOperatorFile", func(t *testing.T) {
		t.Log("CreateOperatorFile", randomData)
		err := client.CreateOperatorFile(context.Background(), owner, repoName, randomData, path, "test-push")
		assert.Nil(t, err)
	})
}

func TestGithubClient_GetOperatorFileSHA(t *testing.T) {
	t.Run("GetOperatorFileSHA", func(t *testing.T) {
		sha, err := client.GetOperatorFileSHA(context.Background(), owner, repoName, path)
		assert.Nil(t, err)
		assert.NotEmpty(t, sha)
	})
	t.Run("Not Exists GetOperatorFileSHA", func(t *testing.T) {
		sha, err := client.GetOperatorFileSHA(context.Background(), owner, repoName, "daaa/operator.yaml")
		assert.Nil(t, err)
		assert.Nil(t, sha)
	})
}

func TestGithubClient_DeleteOperatorFile(t *testing.T) {
	t.Run("DeleteOperatorFile", func(t *testing.T) {
		err := client.DeleteOperatorFile(context.Background(), owner, repoName, path)
		t.Log(err)
	})
}
