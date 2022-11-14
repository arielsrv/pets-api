package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ent"
	"github.com/ent/secret"
	"github.com/internal/clients/gitlab"
	"github.com/internal/infrastructure"
	"github.com/internal/model"
	"github.com/internal/shared"
)

type ISecretService interface {
	SaveSecret(appId int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error)
	GetSecret(secretID int64) (string, error)
}

type SecretService struct {
	client     gitlab.IGitLabClient
	dataAccess *infrastructure.DataAccessService
	appService IAppService
}

func NewSecretService(client gitlab.IGitLabClient, dataAccess *infrastructure.DataAccessService, appService IAppService) *SecretService {
	return &SecretService{client: client, dataAccess: dataAccess, appService: appService}
}

func (s *SecretService) GetSecret(secretID int64) (string, error) {
	result, err := s.dataAccess.GetClient().Secret.Query().
		Where(secret.ID(secretID)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return "", shared.NewError(http.StatusNotFound, fmt.Sprintf("secret with id %d not found", secretID))
		}
		return "", err
	}

	return result.Key, nil
}

func (s *SecretService) SaveSecret(appId int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error) {
	appModel, err := s.appService.GetAppById(appId)
	if err != nil {
		return nil, err
	}

	secretName := fmt.Sprintf("PETS_%s_%s", strings.ToUpper(appModel.Name), strings.ToUpper(secretModel.Key))

	result, err := s.dataAccess.GetClient().Secret.Create().
		SetKey(secretName).
		SetValue(secretModel.Value).
		SetAppID(appModel.ID).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	model := new(model.AppSecretModel)
	model.Key = result.Key
	model.RelativeUrl = fmt.Sprintf("/apps/%d/secrets/%d/snippets", appId, result.ID)

	return model, nil
}
