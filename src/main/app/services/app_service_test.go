package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/src/main/app/clients/gitlab/requests"
	"github.com/src/main/app/clients/gitlab/responses"
	"github.com/src/main/app/config"
	"github.com/src/main/app/infrastructure"
	"github.com/src/main/app/model"
	"github.com/src/main/app/services"

	"github.com/src/main/app/ent"

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

	dbClient := infrastructure.NewDbClient(infrastructure.NewSQLiteClient(t))
	dbClient.Open()
	defer dbClient.Close()

	service := services.NewAppService(client, dbClient)
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

	dbClient := infrastructure.NewDbClient(infrastructure.NewSQLiteClient(t))
	dbClient.Open()
	defer dbClient.Close()

	service := services.NewAppService(client, dbClient)
	actual, err := service.GetGroups()

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func GetGroupsError() ([]responses.GroupResponse, error) {
	return nil, errors.New("src server error")
}

func TestAppService_CreateRepository(t *testing.T) {
	client := new(MockClient)
	client.On("CreateProject").Return(GetCreateProjectResponse())

	dbClient := infrastructure.NewDbClient(infrastructure.NewSQLiteClient(t))
	dbClient.Open()
	defer dbClient.Close()

	dbClient.AppType.Create().SetID(1).SetName("backend").Save(context.Background())

	service := services.NewAppService(client, dbClient)
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

	dbClient := infrastructure.NewDbClient(infrastructure.NewSQLiteClient(t))
	dbClient.Open()
	defer dbClient.Close()

	dbClient.AppType.Create().SetID(1).SetName("backend").Save(context.Background())

	service := services.NewAppService(client, dbClient)
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

func TestAppService_GetAppByName(t *testing.T) {
	client := new(MockClient)
	client.On("GetProject").Return(GetProject())

	dbClient := infrastructure.NewDbClient(infrastructure.NewSQLiteClient(t))
	dbClient.Open()
	defer dbClient.Close()

	dbClient.AppType.Create().SetID(1).SetName("backend").Save(context.Background())
	dbClient.App.Create().SetName("customers-api").SetProjectId(1).SetAppTypeID(1).Save(context.Background())

	service := services.NewAppService(client, dbClient)
	actual, err := service.GetAppByName("customers-api")

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, int64(1), actual.ID)
	assert.Equal(t, fmt.Sprintf("https://oauth2:%s@domain.com/repo_url", config.String("gitlab.token")), actual.URL)
}

func TestAppService_GetAppById(t *testing.T) {
	client := new(MockClient)
	client.On("GetProject").Return(GetProject())

	dbClient := infrastructure.NewDbClient(infrastructure.NewSQLiteClient(t))
	dbClient.Open()
	defer dbClient.Close()

	dbClient.AppType.Create().SetID(1).SetName("backend").Save(context.Background())
	dbClient.App.Create().SetName("customers-api").SetProjectId(1).SetAppTypeID(1).Save(context.Background())

	service := services.NewAppService(client, dbClient)
	actual, err := service.GetAppById(1)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, int64(1), actual.ID)
	assert.Equal(t, fmt.Sprintf("https://oauth2:%s@domain.com/repo_url", config.String("gitlab.token")), actual.URL)
}

func TestAppService_GetApp_NotFoundErr(t *testing.T) {
	client := new(MockClient)
	client.On("GetProject").Return(GetProjectNotFoundErr())

	dbClient := infrastructure.NewDbClient(infrastructure.NewSQLiteClient(t))
	dbClient.Open()
	defer dbClient.Close()

	service := services.NewAppService(client, dbClient)
	actual, err := service.GetAppByName("loyalty-api")

	assert.Error(t, err)
	assert.Nil(t, actual)
	assert.Equal(t, "application with name loyalty-api not found", err.Error())
}

func TestAppService_GetAppTypes(t *testing.T) {
	client := new(MockClient)

	dbClient := infrastructure.NewDbClient(infrastructure.NewSQLiteClient(t))
	dbClient.Open()
	defer dbClient.Close()

	dbClient.AppType.Create().SetID(1).SetName("backend").Save(context.Background())

	service := services.NewAppService(client, dbClient)
	actual, err := service.GetAppTypes()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 1)
	assert.Equal(t, 1, actual[0].ID)
	assert.Equal(t, "backend", actual[0].Name)
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
