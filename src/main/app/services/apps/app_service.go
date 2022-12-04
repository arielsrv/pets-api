package apps

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/src/main/app/infrastructure/secrets"

	"github.com/src/main/app/infrastructure/database"

	"github.com/src/main/app/server"

	"github.com/src/main/app/clients/gitlab"
	"github.com/src/main/app/clients/gitlab/requests"
	"github.com/src/main/app/config"
	"github.com/src/main/app/ent"
	"github.com/src/main/app/ent/app"
	"github.com/src/main/app/model"
)

type IAppService interface {
	GetGroups() ([]model.AppGroup, error)
	CreateApp(createAppModel *model.App) (*model.App, error)
	GetAppTypes() ([]model.AppType, error)
	GetAppByName(appName string) (*model.App, error)
	GetAppByID(appID int64) (*model.App, error)
}

type AppService struct {
	gitLabClient gitlab.IGitLabClient
	dbClient     database.IDbClient
	secretStore  secrets.ISecretStore
}

func NewAppService(gitLabClient gitlab.IGitLabClient, dbClient database.IDbClient, secretStore secrets.ISecretStore) *AppService {
	return &AppService{
		gitLabClient: gitLabClient,
		dbClient:     dbClient,
		secretStore:  secretStore,
	}
}

func (s *AppService) GetAppByName(appName string) (*model.App, error) {
	application, err := s.dbClient.Context().App.Query().
		Where(app.Name(appName)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, server.NewError(http.StatusNotFound, fmt.Sprintf("application with name %s not found", appName))
		}
		return nil, err
	}

	projectResponse, err := s.gitLabClient.GetProject(application.ExternalGitlabProjectID)

	if err != nil {
		return nil, err
	}

	repoURL, err := url.Parse(projectResponse.URL)
	if err != nil {
		return nil, err
	}

	gitlabToken := s.secretStore.GetSecret("GITLAB_TOKEN")
	if gitlabToken.Err != nil {
		return nil, gitlabToken.Err
	}

	secureURL := fmt.Sprintf("%s://oauth2:%s@%s%s",
		repoURL.Scheme,
		gitlabToken.Value,
		repoURL.Host,
		repoURL.Path)

	appModel := new(model.App)
	appModel.ID = application.ID
	appModel.URL = secureURL

	return appModel, nil
}

func (s *AppService) GetAppByID(appID int64) (*model.App, error) {
	application, err := s.dbClient.Context().App.Query().
		Where(app.ID(appID)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, server.NewError(http.StatusNotFound, fmt.Sprintf("application with name %d not found", appID))
		}
		return nil, err
	}

	projectResponse, err := s.gitLabClient.GetProject(application.ExternalGitlabProjectID)

	if err != nil {
		return nil, err
	}

	repoURL, err := url.Parse(projectResponse.URL)
	if err != nil {
		return nil, err
	}

	gitlabToken := s.secretStore.GetSecret("GITLAB_TOKEN")
	if gitlabToken.Err != nil {
		return nil, gitlabToken.Err
	}

	secureURL := fmt.Sprintf("%s://oauth2:%s@%s%s",
		repoURL.Scheme,
		gitlabToken.Value,
		repoURL.Host,
		repoURL.Path)

	appModel := new(model.App)
	appModel.ID = application.ID
	appModel.Name = application.Name
	appModel.URL = secureURL

	return appModel, nil
}

func (s *AppService) GetAppTypes() ([]model.AppType, error) {
	appTypes, err := s.dbClient.Context().AppType.
		Query().
		All(context.Background())

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

func (s *AppService) GetGroups() ([]model.AppGroup, error) {
	groupsResponse, err := s.gitLabClient.GetGroups()
	if err != nil {
		return nil, err
	}

	var groupsDto []model.AppGroup
	for _, groupResponse := range groupsResponse {
		var groupDto model.AppGroup
		groupDto.ID = groupResponse.ID
		groupDto.Name = groupResponse.Path
		groupsDto = append(groupsDto, groupDto)
	}

	return groupsDto, nil
}

func (s *AppService) CreateApp(appModel *model.App) (*model.App, error) {
	alreadyExist, err := s.dbClient.Context().App.Query().
		Where(app.Name(appModel.Name)).
		Exist(context.Background())

	if err != nil {
		return nil, err
	}

	if alreadyExist {
		return nil, server.NewError(http.StatusConflict, fmt.Sprintf("project name %s already exist", appModel.Name))
	}

	createProjectRequest := new(requests.CreateProjectRequest)
	createProjectRequest.Name = fmt.Sprintf("%s%s", config.String("gitlab.prefix"), appModel.Name)
	createProjectRequest.GroupID = appModel.GroupID

	response, err := s.gitLabClient.CreateProject(createProjectRequest)

	if err != nil {
		return nil, err
	}

	application, err := s.dbClient.Context().App.Create().
		SetName(appModel.Name).
		SetExternalGitlabProjectID(response.ID).
		SetAppTypeID(int(appModel.AppTypeID)).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	appModel.ID = application.ID
	appModel.URL = response.URL

	return appModel, err
}
