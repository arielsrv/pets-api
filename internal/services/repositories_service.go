package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ent"
	"github.com/internal/server"

	"github.com/ent/app"
	"github.com/internal/clients"
	"github.com/internal/clients/requests"
	"github.com/internal/infrastructure"
	"github.com/internal/model"
	"github.com/internal/shared"
)

type IRepositoriesService interface {
	GetGroups() ([]model.GroupModel, error)
	CreateRepository(repositoryDto *model.RepositoryModel) (*model.AppModel, error)
	GetAppTypes() ([]model.AppType, error)
	GetApp(appName string) (*model.AppModel, error)
}

type RepositoriesService struct {
	client     clients.IGitLabClient
	dataAccess *infrastructure.DataAccessService
}

func NewRepositoriesService(client clients.IGitLabClient, dataAccess *infrastructure.DataAccessService) *RepositoriesService {
	return &RepositoriesService{client: client, dataAccess: dataAccess}
}

func (s *RepositoriesService) GetApp(appName string) (*model.AppModel, error) {
	app, err := s.dataAccess.GetClient().App.Query().
		Where(app.Name(appName)).
		First(context.Background())

	if err != nil {
		var notFoundErr *ent.NotFoundError
		if errors.As(err, &notFoundErr) {
			return nil, shared.NewError(http.StatusNotFound, fmt.Sprintf("app with name %s not found", appName))
		}
		return nil, err
	}

	projectResponse, err := s.client.GetProject(app.ProjectId)

	if err != nil {
		return nil, err
	}

	repoURL, err := url.Parse(projectResponse.URL)
	if err != nil {
		return nil, err
	}

	secureURL := fmt.Sprintf("%s://oauth2:%s@%s%s",
		repoURL.Scheme,
		server.GetAppConfig().GitLab.Token,
		repoURL.Host,
		repoURL.Path)

	appModel := new(model.AppModel)
	appModel.ID = app.ID
	appModel.URL = secureURL

	return appModel, nil
}

func (s *RepositoriesService) GetAppTypes() ([]model.AppType, error) {
	appTypes, err := s.dataAccess.GetClient().AppType.Query().All(context.Background())

	if err != nil {
		return nil, err
	}

	var appTypesModel []model.AppType
	for _, appType := range appTypes {
		var appTypeModel model.AppType
		appTypeModel.ID = appType.ID
		appTypeModel.Name = appType.Name
		appTypesModel = append(appTypesModel, appTypeModel)
	}

	return appTypesModel, nil
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

func (s *RepositoriesService) CreateRepository(repositoryDto *model.RepositoryModel) (*model.AppModel, error) {
	duplicated, err := s.dataAccess.GetClient().App.Query().
		Where(app.Name(repositoryDto.Name)).
		Exist(context.Background())

	if err != nil {
		return nil, err
	}

	if duplicated {
		return nil, shared.NewError(http.StatusConflict, fmt.Sprintf("duplicated project name %s", repositoryDto.Name))
	}

	createProjectRequest := new(requests.CreateProjectRequest)
	createProjectRequest.Name = repositoryDto.Name
	createProjectRequest.GroupID = repositoryDto.GroupID

	response, err := s.client.CreateProject(createProjectRequest)

	if err != nil {
		return nil, err
	}

	application, err := s.dataAccess.GetClient().App.Create().
		SetName(createProjectRequest.Name).
		SetProjectId(response.ID).
		SetAppTypeID(repositoryDto.AppTypeID).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	repoURL, err := url.Parse(response.URL)
	if err != nil {
		return nil, err
	}

	secureURL := fmt.Sprintf("%s://oauth2:%s@%s%s",
		repoURL.Scheme,
		server.GetAppConfig().GitLab.Token,
		repoURL.Host,
		repoURL.Path)

	appModel := new(model.AppModel)
	appModel.ID = application.ID
	appModel.URL = secureURL

	return appModel, err
}
