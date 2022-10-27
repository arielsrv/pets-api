package services

import (
	"github.com/internal/clients"
	"github.com/internal/clients/requests"
	"github.com/internal/model"
)

type IRepositoriesService interface {
	GetGroups() ([]model.GroupModel, error)
	CreateRepository(repositoryDto *model.RepositoryModel) error
}

type RepositoriesService struct {
	client clients.IGitLabClient
}

func NewRepositoriesService(client clients.IGitLabClient) *RepositoriesService {
	return &RepositoriesService{client: client}
}

func (service *RepositoriesService) GetGroups() ([]model.GroupModel, error) {
	groupsResponse, err := service.client.GetGroups()
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

func (service *RepositoriesService) CreateRepository(repositoryDto *model.RepositoryModel) error {
	request := new(requests.CreateProjectRequest)
	request.Name = repositoryDto.Name
	request.GroupID = repositoryDto.GroupID

	err := service.client.CreateProject(request)
	return err
}
