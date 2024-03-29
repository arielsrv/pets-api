package ensure_test

import (
	"testing"

	"github.com/arielsrv/pets-api/src/main/app/ent/property"
	"github.com/arielsrv/pets-api/src/main/app/helpers/ensure"
	"github.com/stretchr/testify/assert"
)

func TestEnsureNotEmpty(t *testing.T) {
	err := ensure.NotEmpty("", "empty value")
	assert.Error(t, err)
	assert.Equal(t, "empty value", err.Error())
}

func TestEnsureInt(t *testing.T) {
	err := ensure.Int(0, "invalid value")
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}

func TestEnsureInt64(t *testing.T) {
	err := ensure.Int64(int64(0), "invalid value")
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}

func TestEnsureNotEmpty_Ok(t *testing.T) {
	err := ensure.NotEmpty("value", "empty value")
	assert.NoError(t, err)
}

func TestEnsureInt_Ok(t *testing.T) {
	err := ensure.Int(1, "invalid value")
	assert.NoError(t, err)
}

func TestEnsureInt64_Ok(t *testing.T) {
	err := ensure.Int64(int64(1), "invalid value")
	assert.NoError(t, err)
}

func TestEnsureEnum(t *testing.T) {
	err := ensure.Enum(property.Backend, "invalid value")
	assert.NoError(t, err)

	err = ensure.Enum(property.Frontend, "invalid value")
	assert.NoError(t, err)
}
