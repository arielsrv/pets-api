package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/internal/shared"

	"github.com/internal/handlers"
	"github.com/internal/model"
	"github.com/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_ "github.com/mattn/go-sqlite3"
)

type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) GetGroups() ([]model.AppGroupModel, error) {
	args := m.Called()
	return args.Get(0).([]model.AppGroupModel), args.Error(1)
}

func (m *MockAppService) CreateApp(*model.CreateAppModel) (*model.AppModel, error) {
	args := m.Called()
	return args.Get(0).(*model.AppModel), args.Error(1)
}

func (m *MockAppService) GetAppTypes() ([]model.AppType, error) {
	args := m.Called()
	return args.Get(0).([]model.AppType), args.Error(1)
}

func (m *MockAppService) GetApp(string) (*model.AppModel, error) {
	args := m.Called()
	return args.Get(0).(*model.AppModel), args.Error(1)
}

func TestRepositoriesHandler_GetGroups(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetGroups").Return(GetGroups())

	repositoriesHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodGet, "/apps/groups", repositoriesHandler.GetGroups)

	request := httptest.NewRequest(http.MethodGet, "/apps/groups", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "[{\"id\":1,\"name\":\"root/group1\"},{\"id\":2,\"name\":\"root/group2\"}]", string(body))
}

func TestRepositoriesHandler_GetGroups_Err(t *testing.T) {
	repositoriesService := new(MockAppService)
	repositoriesService.On("GetGroups").Return(GetGroupsErr())
	repositoriesHandler := handlers.NewAppHandler(repositoriesService)
	app := server.New()
	app.Add(http.MethodGet, "/repositories/groups", repositoriesHandler.GetGroups)

	request := httptest.NewRequest(http.MethodGet, "/repositories/groups", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"status_code\":500,\"message\":\"internal server error\"}", string(body))
}

func GetGroupsErr() ([]model.AppGroupModel, error) {
	return nil, errors.New("internal server error")
}

func TestRepositoriesHandler_CreateApp(t *testing.T) {
	appService := new(MockAppService)
	appService.On("CreateApp").Return(GetCreateApp())
	repositoriesHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodPost, "/apps", repositoriesHandler.CreateApp)

	request := httptest.
		NewRequest(http.MethodPost, "/apps",
			bytes.NewBufferString("{\"name\":\"my repo\",\"group_id\":1, \"app_type_id\": 1}"))

	request.Header.Add("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"id\":1}", string(body))
}

func GetCreateApp() (*model.AppModel, error) {
	appModel := new(model.AppModel)
	appModel.ID = 1
	return appModel, nil
}

func TestRepositoriesHandler_CreateApp_Err(t *testing.T) {
	repositoriesService := new(MockAppService)
	repositoriesService.On("CreateApp").Return(GetCreateError())
	repositoriesHandler := handlers.NewAppHandler(repositoriesService)
	app := server.New()
	app.Add(http.MethodPost, "/repositories", repositoriesHandler.CreateApp)

	request := httptest.
		NewRequest(http.MethodPost, "/repositories",
			bytes.NewBufferString("{\"name\":\"my repo\",\"group_id\":1, \"app_type_id\": 1}"))

	request.Header.Add("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"status_code\":500,\"message\":\"internal server error\"}", string(body))
}

func GetCreateError() (*model.AppModel, error) {
	return nil, errors.New("internal server error")
}

func TestRepositoriesHandler_CreateApp_BadRequest_Err(t *testing.T) {
	repositoriesService := new(MockAppService)
	repositoriesService.On("CreateApp").Return(shared.NewError(http.StatusBadRequest, "bad request error"))
	repositoriesHandler := handlers.NewAppHandler(repositoriesService)
	app := server.New()
	app.Add(http.MethodPost, "/repositories", repositoriesHandler.CreateApp)

	request := httptest.
		NewRequest(http.MethodPost, "/repositories",
			bytes.NewBufferString("{\"invalid_field\":\"my repo\",\"group_id\":1}"))

	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"status_code\":400,\"message\":\"bad request error\"}", string(body))
}

func GetGroups() ([]model.AppGroupModel, error) {
	var group1 model.AppGroupModel
	group1.ID = 1
	group1.Name = "root/group1"
	var group2 model.AppGroupModel
	group2.ID = 2
	group2.Name = "root/group2"

	var groups []model.AppGroupModel
	groups = append(groups, group1)
	groups = append(groups, group2)

	return groups, nil
}
