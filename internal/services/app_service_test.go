package services_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/internal/clients/gitlab/responses"
	"github.com/internal/config"

	"github.com/internal/clients/gitlab/requests"

	"github.com/ent"

	"github.com/internal/infrastructure"

	"github.com/internal/model"

	"github.com/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetGroups() ([]responses.GroupResponse, error) {
	args := m.Called()
	return args.Get(0).([]responses.GroupResponse), args.Error(1)
}

func (m *MockClient) CreateProject(*requests.CreateProjectRequest) (*responses.CreateProjectResponse, error) {
	args := m.Called()
	return args.Get(0).(*responses.CreateProjectResponse), args.Error(1)
}

func (m *MockClient) GetProject(int64) (*responses.ProjectResponse, error) {
	args := m.Called()
	return args.Get(0).(*responses.ProjectResponse), args.Error(1)
}

func TestAppService_GetGroups(t *testing.T) {
	client := new(MockClient)
	client.On("GetGroups").Return(GetGroups())

	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewAppService(client, dataAccessService)
	actual, err := service.GetGroups()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 2)
	assert.Equal(t, int64(1), actual[0].ID)
	assert.Equal(t, "root/group1", actual[0].Name)
	assert.Equal(t, int64(2), actual[1].ID)
	assert.Equal(t, "root/group2", actual[1].Name)
}

func TestAppService_GetGroups_Err(t *testing.T) {
	client := new(MockClient)
	client.On("GetGroups").Return(GetGroupsError())

	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewAppService(client, dataAccessService)
	actual, err := service.GetGroups()

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func GetGroupsError() ([]responses.GroupResponse, error) {
	return nil, errors.New("internal server error")
}

func TestAppService_CreateRepository(t *testing.T) {
	client := new(MockClient)
	client.On("CreateProject").Return(GetCreateProjectResponse())
	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewAppService(client, dataAccessService)
	appModel := new(model.CreateAppModel)
	appModel.Name = "my project name"
	appModel.AppTypeID = 1
	actual, err := service.CreateApp(appModel)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func TestAppService_CreateApp_Conflict(t *testing.T) {
	client := new(MockClient)
	client.On("CreateProject").Return(GetCreateProjectResponse())
	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewAppService(client, dataAccessService)
	repositoryModel := new(model.CreateAppModel)
	repositoryModel.Name = "users-api"
	repositoryModel.GroupID = 1
	repositoryModel.AppTypeID = 1
	actual, err := service.CreateApp(repositoryModel)
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	actual, err = service.CreateApp(repositoryModel)
	assert.Error(t, err)
	assert.Nil(t, actual)
	assert.Equal(t, "duplicated project name users-api", err.Error())
}

func TestAppService_GetApp(t *testing.T) {
	client := new(MockClient)
	client.On("GetProject").Return(GetProject())

	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewAppService(client, dataAccessService)
	actual, err := service.GetAppByName("customers-api")

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, int64(1), actual.ID)
	assert.Equal(t, fmt.Sprintf("https://oauth2:%s@domain.com/repo_url", config.String("gitlab.token")), actual.URL)
}

func TestAppService_GetApp_NotFoundErr(t *testing.T) {
	client := new(MockClient)
	client.On("GetProject").Return(GetProjectNotFoundErr())

	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewAppService(client, dataAccessService)
	actual, err := service.GetAppByName("loyalty-api")

	assert.Error(t, err)
	assert.Nil(t, actual)
	assert.Equal(t, "application with name loyalty-api not found", err.Error())
}

func TestAppService_GetAppTypes(t *testing.T) {
	client := new(MockClient)
	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewAppService(client, dataAccessService)
	actual, err := service.GetAppTypes()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 1)
	assert.Equal(t, 1, actual[0].ID)
	assert.Equal(t, "backend", actual[0].Name)
}

func TestAppService_GetAppTypes_Err(t *testing.T) {
	client := new(MockClient)
	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	dataAccessService.Close()

	service := services.NewAppService(client, dataAccessService)
	actual, err := service.GetAppTypes()

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func GetProject() (*responses.ProjectResponse, error) {
	projectResponse := new(responses.ProjectResponse)
	projectResponse.ID = 1
	projectResponse.URL = "https://domain.com/repo_url"

	return projectResponse, nil
}

func GetProjectNotFoundErr() (*responses.ProjectResponse, error) {
	var notFoundErr *ent.NotFoundError
	return nil, notFoundErr
}

func GetCreateProjectResponse() (*responses.CreateProjectResponse, error) {
	var createProjectResponse responses.CreateProjectResponse
	createProjectResponse.ID = 1
	createProjectResponse.URL = "https://gitlab.com/repoURL"

	return &createProjectResponse, nil
}

func GetGroups() ([]responses.GroupResponse, error) {
	var group1 responses.GroupResponse
	group1.ID = 1
	group1.Path = "root/group1"
	var group2 responses.GroupResponse
	group2.ID = 2
	group2.Path = "root/group2"

	var groups []responses.GroupResponse
	groups = append(groups, group1)
	groups = append(groups, group2)

	return groups, nil
}
