package secrets_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src/main/app/handlers/secrets"

	"github.com/src/main/app/model"
	"github.com/src/main/app/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSecretService struct {
	mock.Mock
}

func (m *MockSecretService) CreateSecret(int64, *model.CreateSecretRequest) (*model.CreateSecretResponse, error) {
	args := m.Called()
	return args.Get(0).(*model.CreateSecretResponse), args.Error(1)
}

func (m *MockSecretService) GetSecret(int64) (string, error) {
	// TODO implement me
	panic("implement me")
}

type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) GetGroups() ([]model.AppGroupResponse, error) {
	args := m.Called()
	return args.Get(0).([]model.AppGroupResponse), args.Error(1)
}

func (m *MockAppService) CreateApp(*model.CreateAppRequest) (*model.CreateAppResponse, error) {
	args := m.Called()
	return args.Get(0).(*model.CreateAppResponse), args.Error(1)
}

func (m *MockAppService) GetAppTypes() ([]model.AppTypeResponse, error) {
	args := m.Called()
	return args.Get(0).([]model.AppTypeResponse), args.Error(1)
}

func (m *MockAppService) GetAppByName(string) (*model.AppResponse, error) {
	args := m.Called()
	return args.Get(0).(*model.AppResponse), args.Error(1)
}

func (m *MockAppService) GetAppByID(int64) (*model.AppResponse, error) {
	args := m.Called()
	return args.Get(0).(*model.AppResponse), args.Error(1)
}

func TestSecretHandler_CreateSecret(t *testing.T) {
	appService := new(MockAppService)
	secretService := new(MockSecretService)
	secretService.On("CreateSecret").Return(GetNewSecret())

	secretHandler := secrets.NewSecretHandler(appService, secretService)
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

	assert.Equal(t, "{\"key\":\"my_secret_key\",\"snippet_url\":\"/relative_url\"}", string(body))
}

func GetNewSecret() (*model.CreateSecretResponse, error) {
	appSecretModel := new(model.CreateSecretResponse)
	appSecretModel.Key = "my_secret_key"
	appSecretModel.SnippetURL = "/relative_url"

	return appSecretModel, nil
}
