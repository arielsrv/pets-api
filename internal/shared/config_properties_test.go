package shared_test

import (
	"github.com/internal/shared"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetProperty(t *testing.T) {
	actual := shared.GetProperty("gitlab.prefix")
	assert.Equal(t, "ikp_", actual)
}

func TestGetProperty_Err(t *testing.T) {
	actual := shared.GetProperty("missing")
	assert.Equal(t, "", actual)
}

func TestGetIntProperty(t *testing.T) {
	actual := shared.GetIntProperty("database.port")
	assert.Equal(t, 3306, actual)
}

func TestGetIntProperty_Err(t *testing.T) {
	actual := shared.GetIntProperty("missing")
	assert.Equal(t, 0, actual)
}
