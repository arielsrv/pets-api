package files_test

import (
	"github.com/src/main/app/helpers/files"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExist(t *testing.T) {
	actual := files.Exist("files_test.go")
	assert.True(t, actual)
}

func TestExist_NotExist(t *testing.T) {
	actual := files.Exist("files_test2.go")
	assert.False(t, actual)
}

func TestExist_IsDir(t *testing.T) {
	actual := files.Exist("../files")
	assert.False(t, actual)
}
