package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/src/main/model"

	"github.com/src/main/clients/gitlab"
	"github.com/src/main/clients/gitlab/requests"
	"github.com/src/main/config"
	"github.com/src/main/infrastructure"
	"github.com/src/main/shared"

	"github.com/src/main/ent"
	"github.com/src/main/ent/app"
)

type IAppService interface {
	GetGroups() ([]model.AppGroupModel, error)
	CreateApp(createAppModel *model.CreateAppModel) (*model.AppModel, error)
	GetAppTypes() ([]model.AppTypeModel, error)
	GetAppByName(appName string) (*model.AppModel, error)
	GetAppById(appId int64) (*model.AppModel, error)
}

type AppService struct {
	gitLabClient gitlab.IGitLabClient
	dbClient     *infrastructure.DbClient
}

func NewAppService(gitLabClient gitlab.IGitLabClient, dbClient *infrastructure.DbClient) *AppService {
	return &AppService{gitLabClient: gitLabClient, dbClient: dbClient}
}

func (s *AppService) GetAppByName(appName string) (*model.AppModel, error) {
	application, err := s.dbClient.GetClient().App.Query().
		Where(app.Name(appName)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, shared.NewError(http.StatusNotFound, fmt.Sprintf("application with name %s not found", appName))
		}
		return nil, err
	}

	projectResponse, err := s.gitLabClient.GetProject(application.ProjectId)

	if err != nil {
		return nil, err
	}

	repoURL, err := url.Parse(projectResponse.URL)
	if err != nil {
		return nil, err
	}

	gitlabToken := config.String("gitlab.token")

	secureURL := fmt.Sprintf("%s://oauth2:%s@%s%s",
		repoURL.Scheme,
		gitlabToken,
		repoURL.Host,
		repoURL.Path)

	appModel := new(model.AppModel)
	appModel.ID = application.ID
	appModel.URL = secureURL

	return appModel, nil
}

func (s *AppService) GetAppById(appId int64) (*model.AppModel, error) {
	application, err := s.dbClient.GetClient().App.Query().
		Where(app.ID(appId)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, shared.NewError(http.StatusNotFound, fmt.Sprintf("application with name %d not found", appId))
		}
		return nil, err
	}

	projectResponse, err := s.gitLabClient.GetProject(application.ProjectId)

	if err != nil {
		return nil, err
	}

	repoURL, err := url.Parse(projectResponse.URL)
	if err != nil {
		return nil, err
	}

	gitlabToken := config.String("gitlab.token")

	secureURL := fmt.Sprintf("%s://oauth2:%s@%s%s",
		repoURL.Scheme,
		gitlabToken,
		repoURL.Host,
		repoURL.Path)

	appModel := new(model.AppModel)
	appModel.ID = application.ID
	appModel.Name = application.Name
	appModel.URL = secureURL

	return appModel, nil
}

func (s *AppService) GetAppTypes() ([]model.AppTypeModel, error) {
	appTypes, err := s.dbClient.GetClient().AppType.Query().All(context.Background())

	if err != nil {
		return nil, err
	}

	var appTypesModel []model.AppTypeModel
	for _, appType := range appTypes {
		var appTypeModel model.AppTypeModel
		appTypeModel.ID = appType.ID
		appTypeModel.Name = appType.Name
		appTypesModel = append(appTypesModel, appTypeModel)
	}

	return appTypesModel, nil
}

func (s *AppService) GetGroups() ([]model.AppGroupModel, error) {
	groupsResponse, err := s.gitLabClient.GetGroups()
	if err != nil {
		return nil, err
	}

	var groupsDto []model.AppGroupModel
	for _, groupResponse := range groupsResponse {
		var groupDto model.AppGroupModel
		groupDto.ID = groupResponse.ID
		groupDto.Name = groupResponse.Path
		groupsDto = append(groupsDto, groupDto)
	}

	return groupsDto, nil
}

func (s *AppService) CreateApp(createAppModel *model.CreateAppModel) (*model.AppModel, error) {
	duplicated, err := s.dbClient.GetClient().App.Query().
		Where(app.Name(createAppModel.Name)).
		Exist(context.Background())

	if err != nil {
		return nil, err
	}

	if duplicated {
		return nil, shared.NewError(http.StatusConflict, fmt.Sprintf("duplicated project name %s", createAppModel.Name))
	}

	createProjectRequest := new(requests.CreateProjectRequest)
	createProjectRequest.Name = fmt.Sprintf("%s%s", config.String("gitlab.prefix"), createAppModel.Name)
	createProjectRequest.GroupID = createAppModel.GroupID

	response, err := s.gitLabClient.CreateProject(createProjectRequest)

	if err != nil {
		return nil, err
	}

	application, err := s.dbClient.GetClient().App.Create().
		SetName(createAppModel.Name).
		SetProjectId(response.ID).
		SetAppTypeID(int(createAppModel.AppTypeID)).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	appModel := new(model.AppModel)
	appModel.ID = application.ID
	appModel.URL = response.URL

	return appModel, err
}
