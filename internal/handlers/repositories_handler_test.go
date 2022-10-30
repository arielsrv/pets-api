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

type MockRepositoriesService struct {
	mock.Mock
}

func (m *MockRepositoriesService) GetGroups() ([]model.GroupModel, error) {
	args := m.Called()
	return args.Get(0).([]model.GroupModel), args.Error(1)
}

func (m *MockRepositoriesService) CreateRepository(*model.RepositoryModel) (*model.AppModel, error) {
	args := m.Called()
	return args.Get(0).(*model.AppModel), args.Error(1)
}

func (m *MockRepositoriesService) GetAppTypes() ([]model.AppType, error) {
	args := m.Called()
	return args.Get(0).([]model.AppType), args.Error(1)
}

func (m *MockRepositoriesService) GetApp(string) (*model.AppModel, error) {
	args := m.Called()
	return args.Get(0).(*model.AppModel), args.Error(1)
}

func TestRepositoriesHandler_GetGroups(t *testing.T) {
	repositoriesService := new(MockRepositoriesService)
	repositoriesService.On("GetGroups").Return(GetGroups())

	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)
	app := server.New()
	app.Add(http.MethodGet, "/repositories/groups", repositoriesHandler.GetGroups)

	request := httptest.NewRequest(http.MethodGet, "/repositories/groups", nil)
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
	repositoriesService := new(MockRepositoriesService)
	repositoriesService.On("GetGroups").Return(GetGroupsErr())
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)
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

func GetGroupsErr() ([]model.GroupModel, error) {
	return nil, errors.New("internal server error")
}

func TestRepositoriesHandler_CreateRepository(t *testing.T) {
	repositoriesService := new(MockRepositoriesService)
	repositoriesService.On("CreateRepository").Return(GetCreateRepository())
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)
	app := server.New()
	app.Add(http.MethodPost, "/repositories", repositoriesHandler.CreateRepository)

	request := httptest.
		NewRequest(http.MethodPost, "/repositories",
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

func GetCreateRepository() (*model.AppModel, error) {
	appModel := new(model.AppModel)
	appModel.ID = 1
	return appModel, nil
}

func TestRepositoriesHandler_CreateRepository_Err(t *testing.T) {
	repositoriesService := new(MockRepositoriesService)
	repositoriesService.On("CreateRepository").Return(GetCreateError())
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)
	app := server.New()
	app.Add(http.MethodPost, "/repositories", repositoriesHandler.CreateRepository)

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

func TestRepositoriesHandler_CreateRepository_BadRequest_Err(t *testing.T) {
	repositoriesService := new(MockRepositoriesService)
	repositoriesService.On("CreateRepository").Return(shared.NewError(http.StatusBadRequest, "bad request error"))
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)
	app := server.New()
	app.Add(http.MethodPost, "/repositories", repositoriesHandler.CreateRepository)

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

func GetGroups() ([]model.GroupModel, error) {
	var group1 model.GroupModel
	group1.ID = 1
	group1.Name = "root/group1"
	var group2 model.GroupModel
	group2.ID = 2
	group2.Name = "root/group2"

	var groups []model.GroupModel
	groups = append(groups, group1)
	groups = append(groups, group2)

	return groups, nil
}
