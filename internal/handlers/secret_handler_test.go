package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/internal/handlers"
	"github.com/internal/model"
	"github.com/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSecretService struct {
	mock.Mock
}

func (m *MockSecretService) SaveSecret(int64, *model.CreateAppSecretModel) (*model.AppSecretModel, error) {
	args := m.Called()
	return args.Get(0).(*model.AppSecretModel), args.Error(1)
}

func (m *MockSecretService) GetSecret(int64) (string, error) {
	// TODO implement me
	panic("implement me")
}

func TestSecretHandler_CreateSecret(t *testing.T) {
	appService := new(MockAppService)
	secretService := new(MockSecretService)
	secretService.On("SaveSecret").Return(GetNewSecret())

	secretHandler := handlers.NewSecretHandler(appService, secretService)
	app := server.New()
	app.Add(http.MethodPost, "/apps/:appId/secrets", secretHandler.CreateSecret)

	request := httptest.
		NewRequest(http.MethodPost, "/apps/1/secrets",
			bytes.NewBufferString("{\"key\":\"my_secret_key\",\"value\":\"my_secret_value\"}"))

	request.Header.Add("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"key\":\"my_secret_key\",\"url\":\"http://example.com/relative_url\"}", string(body))
}

func GetNewSecret() (*model.AppSecretModel, error) {
	model := new(model.AppSecretModel)
	model.Key = "my_secret_key"
	model.RelativeUrl = "/relative_url"

	return model, nil
}
