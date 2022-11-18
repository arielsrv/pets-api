package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/src/main/infrastructure"
	"github.com/src/main/model"
	"github.com/src/main/shared"

	"github.com/ent"
	"github.com/ent/secret"
)

type ISecretService interface {
	SaveSecret(appId int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error)
	GetSecret(secretID int64) (string, error)
}

type SecretService struct {
	dbClient   *infrastructure.DbClient
	appService IAppService
}

func NewSecretService(dbClient *infrastructure.DbClient, appService IAppService) *SecretService {
	return &SecretService{dbClient: dbClient, appService: appService}
}

func (s *SecretService) GetSecret(secretID int64) (string, error) {
	result, err := s.dbClient.GetClient().Secret.Query().
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

	alreadyExist, err := s.dbClient.GetClient().Secret.
		Query().
		Where(secret.Key(secretName)).
		Exist(context.Background())

	if err != nil {
		return nil, err
	}

	if alreadyExist {
		return nil, shared.NewError(http.StatusConflict, fmt.Sprintf("Secret %s for app already exist. ", secretModel.Key))
	}

	result, err := s.dbClient.GetClient().Secret.Create().
		SetKey(secretName).
		SetValue(secretModel.Value).
		SetAppID(appModel.ID).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	model := new(model.AppSecretModel)
	model.OriginalKey = secretModel.Key
	model.Key = result.Key
	model.SnippetUrl = fmt.Sprintf("/apps/%d/secrets/%d/snippets", appId, result.ID)

	return model, nil
}
