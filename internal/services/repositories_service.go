package services

import (
	"github.com/internal/clients"
	"github.com/internal/model"
	"github.com/internal/model/requests"
)

type RepositoriesService struct {
	client *clients.GitLabClient
}

func NewRepositoriesService(client *clients.GitLabClient) *RepositoriesService {
	return &RepositoriesService{client: client}
}

func (service *RepositoriesService) GetGroups() ([]model.GroupDto, error) {
	groupsResponse, err := service.client.GetGroups()
	if err != nil {
		return nil, err
	}
	var groupsDto []model.GroupDto
	for _, groupResponse := range groupsResponse {
		var groupDto model.GroupDto
		groupDto.ID = groupResponse.ID
		groupDto.Name = groupResponse.Path
		groupsDto = append(groupsDto, groupDto)
	}
	return groupsDto, nil
}

func (service *RepositoriesService) CreateRepository(repositoryDto *model.RepositoryDto) error {
	request := new(requests.CreateProjectRequest)
	request.Name = repositoryDto.Name
	request.GroupID = repositoryDto.GroupID

	err := service.client.CreateProject(request)
	return err
}
