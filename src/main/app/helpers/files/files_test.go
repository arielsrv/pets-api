package files_test

import (
	"testing"

	"github.com/arielsrv/pets-api/src/main/app/helpers/files"
	"github.com/stretchr/testify/assert"
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
