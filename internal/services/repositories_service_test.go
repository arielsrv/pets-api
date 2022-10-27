package services_test

import (
	"github.com/internal/model"
	"testing"

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

func (m *MockClient) CreateProject(*requests.CreateProjectRequest) error {
	args := m.Called()
	return args.Error(0)
}

func TestRepositoriesService_GetGroups(t *testing.T) {
	client := new(MockClient)
	client.On("GetGroups").Return(GetGroups())
	service := services.NewRepositoriesService(client)
	actual, err := service.GetGroups()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 2)
	assert.Equal(t, int64(1), actual[0].ID)
	assert.Equal(t, "root/group1", actual[0].Name)
	assert.Equal(t, int64(2), actual[1].ID)
	assert.Equal(t, "root/group2", actual[1].Name)
}

func TestRepositoriesService_CreateRepository(t *testing.T) {
	client := new(MockClient)
	client.On("CreateProject").Return(nil)
	service := services.NewRepositoriesService(client)
	repositoryModel := new(model.RepositoryModel)
	err := service.CreateRepository(repositoryModel)

	assert.NoError(t, err)
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
