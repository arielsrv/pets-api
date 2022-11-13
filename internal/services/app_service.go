package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ent/secret"

	"github.com/internal/config"

	"github.com/internal/clients/gitlab"
	"github.com/internal/clients/gitlab/requests"

	"github.com/ent"
	"github.com/ent/app"
	"github.com/internal/infrastructure"
	"github.com/internal/model"
	"github.com/internal/shared"
)

type IAppService interface {
	GetGroups() ([]model.AppGroupModel, error)
	CreateApp(repositoryDto *model.CreateAppModel) (*model.AppModel, error)
	GetAppTypes() ([]model.AppType, error)
	GetAppByName(appName string) (*model.AppModel, error)
	GetAppById(appId int64) (*model.AppModel, error)
	SaveSecret(appId int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error)
	GetSecret(secretID int64) (string, string, error)
}

type AppService struct {
	client     gitlab.IGitLabClient
	dataAccess *infrastructure.DataAccessService
}

func (s *AppService) GetSecret(secretID int64) (string, string, error) {
	result, err := s.dataAccess.GetClient().Secret.Query().
		Where(secret.ID(secretID)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return "", "", shared.NewError(http.StatusNotFound, fmt.Sprintf("secret with id %d not found", secretID))
		}
		return "", "", err
	}

	app, err := result.QueryApp().Only(context.Background())
	if err != nil {
		return "", "", err
	}

	return result.Key, app.Name, nil
}

func (s *AppService) SaveSecret(appId int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error) {
	result, err := s.dataAccess.GetClient().Secret.Create().
		SetKey(secretModel.Key).
		SetValue(secretModel.Value).
		SetAppID(appId).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	model := new(model.AppSecretModel)
	model.Key = secretModel.Key
	model.RelativeUrl = fmt.Sprintf("/apps/%d/secrets/%d/snippets", appId, result.ID)

	return model, nil
}

func NewAppService(client gitlab.IGitLabClient, dataAccess *infrastructure.DataAccessService) *AppService {
	return &AppService{client: client, dataAccess: dataAccess}
}

func (s *AppService) GetAppByName(appName string) (*model.AppModel, error) {
	application, err := s.dataAccess.GetClient().App.Query().
		Where(app.Name(appName)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, shared.NewError(http.StatusNotFound, fmt.Sprintf("application with name %s not found", appName))
		}
		return nil, err
	}

	projectResponse, err := s.client.GetProject(application.ProjectId)

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
	application, err := s.dataAccess.GetClient().App.Query().
		Where(app.ID(appId)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, shared.NewError(http.StatusNotFound, fmt.Sprintf("application with name %d not found", appId))
		}
		return nil, err
	}

	projectResponse, err := s.client.GetProject(application.ProjectId)

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

func (s *AppService) GetAppTypes() ([]model.AppType, error) {
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

func (s *AppService) GetGroups() ([]model.AppGroupModel, error) {
	groupsResponse, err := s.client.GetGroups()
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

func (s *AppService) CreateApp(repositoryDto *model.CreateAppModel) (*model.AppModel, error) {
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
	createProjectRequest.Name = fmt.Sprintf("%s%s", config.String("gitlab.prefix"), repositoryDto.Name)
	createProjectRequest.GroupID = repositoryDto.GroupID

	response, err := s.client.CreateProject(createProjectRequest)

	if err != nil {
		return nil, err
	}

	application, err := s.dataAccess.GetClient().App.Create().
		SetName(repositoryDto.Name).
		SetProjectId(response.ID).
		SetAppTypeID(repositoryDto.AppTypeID).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	appModel := new(model.AppModel)
	appModel.ID = application.ID
	appModel.URL = response.URL

	return appModel, err
}
