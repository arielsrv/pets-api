package config_test

import (
	"testing"

	"github.com/src/main/app/config"

	"github.com/stretchr/testify/assert"
)

func TestGetProperty(t *testing.T) {
	actual := config.String("gitlab.prefix")
	assert.Equal(t, "ikp_", actual)
}

func TestGetProperty_Err(t *testing.T) {
	actual := config.String("missing")
	assert.Equal(t, "", actual)
}

func TestGetIntProperty(t *testing.T) {
	actual := config.Int("rest.pool.default.pool.size")
	assert.Equal(t, 20, actual)
}

func TestGetIntProperty_Err(t *testing.T) {
	actual := config.Int("missing")
	assert.Equal(t, 0, actual)
}
