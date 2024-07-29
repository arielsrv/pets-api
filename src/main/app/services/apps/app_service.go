package apps

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arielsrv/pets-api/src/main/app/clients/gitlab"
	"github.com/arielsrv/pets-api/src/main/app/config"
	"github.com/arielsrv/pets-api/src/main/app/ent"
	"github.com/arielsrv/pets-api/src/main/app/ent/app"
	"github.com/arielsrv/pets-api/src/main/app/infrastructure/database"
	"github.com/arielsrv/pets-api/src/main/app/infrastructure/secrets"
	"github.com/arielsrv/pets-api/src/main/app/model"
	"github.com/arielsrv/pets-api/src/main/app/server"
)

type IAppService interface {
	GetGroups() ([]model.AppGroupResponse, error)
	CreateApp(createAppModel *model.CreateAppRequest) (*model.CreateAppResponse, error)
	GetAppTypes() ([]model.AppTypeResponse, error)
	GetAppByName(appName string) (*model.AppResponse, error)
	GetAppByID(appID int64) (*model.AppResponse, error)
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

func (s *AppService) GetAppByName(appName string) (*model.AppResponse, error) {
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

	appResponse := new(model.AppResponse)
	appResponse.ID = application.ID
	appResponse.URL = secureURL

	return appResponse, nil
}

func (s *AppService) GetAppByID(appID int64) (*model.AppResponse, error) {
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

	appModel := new(model.AppResponse)
	appModel.ID = application.ID
	appModel.Name = application.Name
	appModel.URL = secureURL

	return appModel, nil
}

func (s *AppService) GetAppTypes() ([]model.AppTypeResponse, error) {
	appTypes, err := s.dbClient.Context().AppType.
		Query().
		All(context.Background())

	if err != nil {
		return nil, err
	}

	appTypesModel := make([]model.AppTypeResponse, 0, len(appTypes))
	for _, appType := range appTypes {
		var appTypeModel model.AppTypeResponse
		appTypeModel.ID = appType.ID
		appTypeModel.Name = appType.Name
		appTypesModel = append(appTypesModel, appTypeModel)
	}

	return appTypesModel, nil
}

func (s *AppService) GetGroups() ([]model.AppGroupResponse, error) {
	groupsResponse, err := s.gitLabClient.GetGroups()
	if err != nil {
		return nil, err
	}

	groupsDto := make([]model.AppGroupResponse, 0, len(groupsResponse))
	for _, groupResponse := range groupsResponse {
		var groupDto model.AppGroupResponse
		groupDto.ID = groupResponse.ID
		groupDto.Name = groupResponse.Path
		groupsDto = append(groupsDto, groupDto)
	}

	return groupsDto, nil
}

func (s *AppService) CreateApp(appModel *model.CreateAppRequest) (*model.CreateAppResponse, error) {
	alreadyExist, err := s.dbClient.Context().App.Query().
		Where(app.Name(appModel.Name)).
		Exist(context.Background())

	if err != nil {
		return nil, err
	}

	if alreadyExist {
		return nil, server.NewError(http.StatusConflict, fmt.Sprintf("project name %s already exist", appModel.Name))
	}

	createProjectRequest := new(gitlab.CreateProjectRequest)
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

	createAppResponse := new(model.CreateAppResponse)
	createAppResponse.ID = application.ID
	createAppResponse.URL = response.URL

	return createAppResponse, err
}
