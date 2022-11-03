package config_test

import (
	"testing"

	"github.com/internal/config"
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
	actual := config.Int("database.port")
	assert.Equal(t, 3306, actual)
}

func TestGetIntProperty_Err(t *testing.T) {
	actual := config.Int("missing")
	assert.Equal(t, 0, actual)
}
