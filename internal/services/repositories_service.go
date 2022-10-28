package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ent/app"
	"github.com/internal/clients"
	"github.com/internal/clients/requests"
	"github.com/internal/infrastructure"
	"github.com/internal/model"
	"github.com/internal/shared"
)

type IRepositoriesService interface {
	GetGroups() ([]model.GroupModel, error)
	CreateRepository(repositoryDto *model.RepositoryModel) (int64, error)
}

type RepositoriesService struct {
	client     clients.IGitLabClient
	dataAccess *infrastructure.DataAccessService
}

func NewRepositoriesService(client clients.IGitLabClient, dataAccess *infrastructure.DataAccessService) *RepositoriesService {
	return &RepositoriesService{client: client, dataAccess: dataAccess}
}

func (s *RepositoriesService) GetGroups() ([]model.GroupModel, error) {
	groupsResponse, err := s.client.GetGroups()
	if err != nil {
		return nil, err
	}
	var groupsDto []model.GroupModel
	for _, groupResponse := range groupsResponse {
		var groupDto model.GroupModel
		groupDto.ID = groupResponse.ID
		groupDto.Name = groupResponse.Path
		groupsDto = append(groupsDto, groupDto)
	}
	return groupsDto, nil
}

func (s *RepositoriesService) CreateRepository(repositoryDto *model.RepositoryModel) (int64, error) {
	duplicated, err := s.dataAccess.GetClient().App.Query().
		Where(app.Name(repositoryDto.Name)).
		Exist(context.Background())

	if err != nil {
		return 0, err
	}

	if duplicated {
		return 0, shared.NewError(http.StatusConflict, fmt.Sprintf("duplicated project name %s", repositoryDto.Name))
	}

	createProjectRequest := new(requests.CreateProjectRequest)
	createProjectRequest.Name = repositoryDto.Name
	createProjectRequest.GroupID = repositoryDto.GroupID

	response, err := s.client.CreateProject(createProjectRequest)

	if err != nil {
		return 0, err
	}

	application, err := s.dataAccess.GetClient().App.Create().
		SetName(createProjectRequest.Name).
		SetProjectId(response.ID).
		Save(context.Background())

	if err != nil {
		return 0, err
	}

	return application.ID, err
}
