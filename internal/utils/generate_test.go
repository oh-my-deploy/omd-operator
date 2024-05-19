package utils_test

import (
	"github.com/oh-my-deploy/omd-operator/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomStringLists(t *testing.T) {
	result := utils.RandomStringLists(2)
	assert.Len(t, result, 2)
}
