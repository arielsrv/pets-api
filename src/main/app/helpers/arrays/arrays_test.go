package arrays_test

import (
	"testing"

	"github.com/src/main/app/helpers/arrays"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	var values []int
	values = append(values, 1, 2)

	actual := arrays.Contains(values, 1)
	assert.True(t, actual)
}

func TestContains_NotFound(t *testing.T) {
	var values []int
	values = append(values, 1, 2)

	actual := arrays.Contains(values, 3)
	assert.False(t, actual)
}
