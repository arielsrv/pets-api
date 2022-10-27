package handlers_test

import (
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

type MockRepositoriesService struct {
	mock.Mock
}

func (m *MockRepositoriesService) GetGroups() ([]model.GroupModel, error) {
	args := m.Called()
	return args.Get(0).([]model.GroupModel), args.Error(1)
}

func (m *MockRepositoriesService) CreateRepository(repositoryDto *model.RepositoryModel) error {
	// TODO implement me
	panic("implement me")
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
