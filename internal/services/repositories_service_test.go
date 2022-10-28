package services_test

import (
	"errors"
	"testing"

	"github.com/internal/infrastructure"

	"github.com/internal/model"

	"github.com/internal/clients/requests"
	"github.com/internal/clients/responses"
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

func TestRepositoriesService_GetGroups(t *testing.T) {
	client := new(MockClient)
	client.On("GetGroups").Return(GetGroups())

	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewRepositoriesService(client, dataAccessService)
	actual, err := service.GetGroups()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 2)
	assert.Equal(t, int64(1), actual[0].ID)
	assert.Equal(t, "root/group1", actual[0].Name)
	assert.Equal(t, int64(2), actual[1].ID)
	assert.Equal(t, "root/group2", actual[1].Name)
}

func TestRepositoriesService_GetGroups_Err(t *testing.T) {
	client := new(MockClient)
	client.On("GetGroups").Return(GetGroupsError())

	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewRepositoriesService(client, dataAccessService)
	actual, err := service.GetGroups()

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func GetGroupsError() ([]responses.GroupResponse, error) {
	return nil, errors.New("internal server error")
}

func TestRepositoriesService_CreateRepository(t *testing.T) {
	client := new(MockClient)
	client.On("CreateProject").Return(GetCreateProjectResponse())
	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewRepositoriesService(client, dataAccessService)
	repositoryModel := new(model.RepositoryModel)
	actual, err := service.CreateRepository(repositoryModel)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func TestRepositoriesService_CreateRepository_Conflict(t *testing.T) {
	client := new(MockClient)
	client.On("CreateProject").Return(GetCreateProjectResponse())
	dataAccessService := infrastructure.NewDataAccessService()
	dataAccessService.Test(t)
	defer dataAccessService.Close()

	service := services.NewRepositoriesService(client, dataAccessService)
	repositoryModel := new(model.RepositoryModel)
	repositoryModel.Name = "project name"
	actual, err := service.CreateRepository(repositoryModel)
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	actual, err = service.CreateRepository(repositoryModel)
	assert.Error(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, "duplicated project name project name", err.Error())
}

func GetCreateProjectResponse() (*responses.CreateProjectResponse, error) {
	var createProjectResponse responses.CreateProjectResponse
	createProjectResponse.ID = 1

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
