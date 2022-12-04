package apps_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src/main/app/handlers/apps"

	"github.com/src/main/app/model"

	"github.com/src/main/app/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_ "github.com/mattn/go-sqlite3"
)

type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) GetGroups() ([]model.AppGroup, error) {
	args := m.Called()
	return args.Get(0).([]model.AppGroup), args.Error(1)
}

func (m *MockAppService) CreateApp(*model.App) (*model.App, error) {
	args := m.Called()
	return args.Get(0).(*model.App), args.Error(1)
}

func (m *MockAppService) GetAppTypes() ([]model.AppType, error) {
	args := m.Called()
	return args.Get(0).([]model.AppType), args.Error(1)
}

func (m *MockAppService) GetAppByName(string) (*model.App, error) {
	args := m.Called()
	return args.Get(0).(*model.App), args.Error(1)
}

func (m *MockAppService) GetAppByID(int64) (*model.App, error) {
	args := m.Called()
	return args.Get(0).(*model.App), args.Error(1)
}

func TestAppHandler_GetGroups(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetGroups").Return(GetGroups())

	appHandler := apps.NewAppHandler(appService)
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

	appHandler := apps.NewAppHandler(appService)
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

	appHandler := apps.NewAppHandler(appService)
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

	appHandler := apps.NewAppHandler(appService)
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

func GetAppNotFound() (*model.App, error) {
	return nil, server.NewError(http.StatusNotFound, "app with name customer-api not found")
}

func GetApp() (*model.App, error) {
	appModel := new(model.App)
	appModel.ID = 1
	appModel.URL = "repo_url"
	return appModel, nil
}

func TestAppHandler_GetGroups_Err(t *testing.T) {
	appService := new(MockAppService)
	appService.On("GetGroups").Return(GetGroupsErr())
	repositoriesHandler := apps.NewAppHandler(appService)
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

func GetGroupsErr() ([]model.AppGroup, error) {
	return nil, errors.New("src server error")
}

func TestRepositoriesHandler_CreateApp(t *testing.T) {
	appService := new(MockAppService)
	appService.On("CreateApp").Return(GetCreateApp())
	appHandler := apps.NewAppHandler(appService)
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

func GetCreateApp() (*model.App, error) {
	appModel := new(model.App)
	appModel.ID = 1
	return appModel, nil
}

func TestRepositoriesHandler_CreateApp_Err(t *testing.T) {
	appService := new(MockAppService)
	appService.On("CreateApp").Return(GetCreateError())
	appHandler := apps.NewAppHandler(appService)
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

func GetCreateError() (*model.App, error) {
	return nil, errors.New("src server error")
}

func TestRepositoriesHandler_CreateApp_BadRequest_Err(t *testing.T) {
	appService := new(MockAppService)
	appService.On("CreateApp").Return(server.NewError(http.StatusBadRequest, "bad request error"))
	appHandler := apps.NewAppHandler(appService)
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

	appHandler := apps.NewAppHandler(appService)
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

	appHandler := apps.NewAppHandler(appService)
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

func GetAppTypesErr() ([]model.AppType, error) {
	return nil, errors.New("src server error")
}

func GetAppTypes() ([]model.AppType, error) {
	var appType1 model.AppType
	appType1.ID = 1
	appType1.Name = "backend"

	var appType2 model.AppType
	appType2.ID = 2
	appType2.Name = "frontend"

	var appsType []model.AppType
	appsType = append(appsType, appType1)
	appsType = append(appsType, appType2)

	return appsType, nil
}

func GetGroups() ([]model.AppGroup, error) {
	var group1 model.AppGroup
	group1.ID = 1
	group1.Name = "root/group1"
	var group2 model.AppGroup
	group2.ID = 2
	group2.Name = "root/group2"

	var groups []model.AppGroup
	groups = append(groups, group1)
	groups = append(groups, group2)

	return groups, nil
}
