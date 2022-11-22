package config_test

import (
	"fmt"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/src/main/app/config"
	"github.com/stretchr/testify/assert"
)

func TestProvideRestClients(t *testing.T) {
	err := setUp("rest_factory_test.yml")
	assert.NoError(t, err)

	restClientFactory := config.ProvideRestClients()
	assert.Greater(t, len(restClientFactory.GetClients()), 1)

	googleClient := restClientFactory.Get("google")
	assert.NotNil(t, googleClient)
	assert.Equal(t, time.Millisecond*1000*2, googleClient.Timeout)
	assert.Equal(t, time.Millisecond*1000*5, googleClient.ConnectTimeout)
	assert.Equal(t, 20, googleClient.CustomPool.MaxIdleConnsPerHost)

	amazonClient := restClientFactory.Get("amazon")
	assert.NotNil(t, amazonClient)
	assert.Equal(t, time.Millisecond*1000*2, amazonClient.Timeout)
	assert.Equal(t, time.Millisecond*1000*5, amazonClient.ConnectTimeout)
	assert.Equal(t, 20, amazonClient.CustomPool.MaxIdleConnsPerHost)

	assert.True(t, amazonClient == googleClient)
}

func TestProvideRestClients_NotReusedPool(t *testing.T) {
	err := setUp("rest_factory_test.yml")
	assert.NoError(t, err)

	restClientFactory := config.ProvideRestClients()

	first := restClientFactory.Get("first")
	assert.NotNil(t, first)
	assert.Equal(t, time.Millisecond*1000*2, first.Timeout)
	assert.Equal(t, time.Millisecond*1000*5, first.ConnectTimeout)
	assert.Equal(t, 20, first.CustomPool.MaxIdleConnsPerHost)

	second := restClientFactory.Get("second")
	assert.NotNil(t, second)
	assert.Equal(t, time.Millisecond*1000*2, second.Timeout)
	assert.Equal(t, time.Millisecond*1000*5, second.ConnectTimeout)
	assert.Equal(t, 20, second.CustomPool.MaxIdleConnsPerHost)

	assert.True(t, second != first)
}

func setUp(file string) error {
	_, caller, _, _ := runtime.Caller(0)
	err := archaius.AddFile(fmt.Sprintf("%s/%s", path.Dir(caller), file))
	return err
}
