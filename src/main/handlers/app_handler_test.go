package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src/main/handlers"
	model2 "github.com/src/main/model"
	"github.com/src/main/server"
	"github.com/src/main/shared"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_ "github.com/mattn/go-sqlite3"
)

type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) GetGroups() ([]model2.AppGroupModel, error) {
	args := m.Called()
	return args.Get(0).([]model2.AppGroupModel), args.Error(1)
}

func (m *MockAppService) CreateApp(*model2.CreateAppModel) (*model2.AppModel, error) {
	args := m.Called()
	return args.Get(0).(*model2.AppModel), args.Error(1)
}

func (m *MockAppService) GetAppTypes() ([]model2.AppTypeModel, error) {
	args := m.Called()
	return args.Get(0).([]model2.AppTypeModel), args.Error(1)
}

func (m *MockAppService) GetAppByName(string) (*model2.AppModel, error) {
	args := m.Called()
	return args.Get(0).(*model2.AppModel), args.Error(1)
}

func (m *MockAppService) GetAppById(int64) (*model2.AppModel, error) {
	args := m.Called()
	return args.Get(0).(*model2.AppModel), args.Error(1)
}

func TestAppHandler_GetGroups(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetGroups").Return(GetGroups())

	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodGet, "/apps/groups", appHandler.GetGroups)

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

func TestAppHandler_GetApp(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetAppByName").Return(GetApp())

	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodGet, "/apps", appHandler.GetApp)

	request := httptest.NewRequest(http.MethodGet, "/apps?app_name=go", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"id\":1,\"url\":\"repo_url\"}", string(body))
}

func TestAppHandler_GetApp_BadRequestErr(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetAppByName").Return(GetApp())

	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodGet, "/apps", appHandler.GetApp)

	request := httptest.NewRequest(http.MethodGet, "/apps", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"status_code\":400,\"message\":\"bad request error, missing app_name\"}", string(body))
}

func TestAppHandler_GetApp_NotFoundErr(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetAppByName").Return(GetAppNotFound())

	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodGet, "/apps", appHandler.GetApp)

	request := httptest.NewRequest(http.MethodGet, "/apps?app_name=customers-api", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"status_code\":404,\"message\":\"app with name customer-api not found\"}", string(body))
}

func GetAppNotFound() (*model2.AppModel, error) {
	return nil, shared.NewError(http.StatusNotFound, "app with name customer-api not found")
}

func GetApp() (*model2.AppModel, error) {
	appModel := new(model2.AppModel)
	appModel.ID = 1
	appModel.URL = "repo_url"
	return appModel, nil
}

func TestAppHandler_GetGroups_Err(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetGroups").Return(GetGroupsErr())
	repositoriesHandler := handlers.NewAppHandler(appService)
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

	assert.Equal(t, "{\"status_code\":500,\"message\":\"src server error\"}", string(body))
}

func GetGroupsErr() ([]model2.AppGroupModel, error) {
	return nil, errors.New("src server error")
}

func TestRepositoriesHandler_CreateApp(t *testing.T) {
	appService := new(MockAppService)
	appService.On("CreateApp").Return(GetCreateApp())
	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodPost, "/apps", appHandler.CreateApp)

	request := httptest.
		NewRequest(http.MethodPost, "/apps",
			bytes.NewBufferString("{\"name\":\"my repo\",\"group_id\":1, \"app_type_id\": 1}"))

	request.Header.Add("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"id\":1}", string(body))
}

func GetCreateApp() (*model2.AppModel, error) {
	appModel := new(model2.AppModel)
	appModel.ID = 1
	return appModel, nil
}

func TestRepositoriesHandler_CreateApp_Err(t *testing.T) {
	appService := new(MockAppService)
	appService.On("CreateApp").Return(GetCreateError())
	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodPost, "/repositories", appHandler.CreateApp)

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

	assert.Equal(t, "{\"status_code\":500,\"message\":\"src server error\"}", string(body))
}

func GetCreateError() (*model2.AppModel, error) {
	return nil, errors.New("src server error")
}

func TestRepositoriesHandler_CreateApp_BadRequest_Err(t *testing.T) {
	appService := new(MockAppService)
	appService.On("CreateApp").Return(shared.NewError(http.StatusBadRequest, "bad request error"))
	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodPost, "/repositories", appHandler.CreateApp)

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

func TestAppHandler_GetAppTypes(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetAppTypes").Return(GetAppTypes())

	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodGet, "/apps/types", appHandler.GetAppTypes)

	request := httptest.NewRequest(http.MethodGet, "/apps/types", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "[{\"id\":1,\"name\":\"backend\"},{\"id\":2,\"name\":\"frontend\"}]", string(body))
}

func TestAppHandler_GetAppTypes_Err(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetAppTypes").Return(GetAppTypesErr())

	appHandler := handlers.NewAppHandler(appService)
	app := server.New()
	app.Add(http.MethodGet, "/apps/types", appHandler.GetAppTypes)

	request := httptest.NewRequest(http.MethodGet, "/apps/types", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"status_code\":500,\"message\":\"src server error\"}", string(body))
}

func GetAppTypesErr() ([]model2.AppTypeModel, error) {
	return nil, errors.New("src server error")
}

func GetAppTypes() ([]model2.AppTypeModel, error) {
	var appType1 model2.AppTypeModel
	appType1.ID = 1
	appType1.Name = "backend"

	var appType2 model2.AppTypeModel
	appType2.ID = 2
	appType2.Name = "frontend"

	var appsType []model2.AppTypeModel
	appsType = append(appsType, appType1)
	appsType = append(appsType, appType2)

	return appsType, nil
}

func GetGroups() ([]model2.AppGroupModel, error) {
	var group1 model2.AppGroupModel
	group1.ID = 1
	group1.Name = "root/group1"
	var group2 model2.AppGroupModel
	group2.ID = 2
	group2.Name = "root/group2"

	var groups []model2.AppGroupModel
	groups = append(groups, group1)
	groups = append(groups, group2)

	return groups, nil
}
