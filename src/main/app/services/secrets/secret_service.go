package secrets

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/src/main/app/services/apps"

	"github.com/src/main/app/infrastructure/database"
	"github.com/src/main/app/server"

	"github.com/src/main/app/ent"
	"github.com/src/main/app/ent/secret"
	"github.com/src/main/app/model"
)

type ISecretService interface {
	CreateSecret(appID int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error)
	GetSecret(secretID int64) (string, error)
}

type SecretService struct {
	dbClient   database.IDbClient
	appService apps.IAppService
}

func NewSecretService(dbClient database.IDbClient, appService apps.IAppService) *SecretService {
	return &SecretService{
		dbClient:   dbClient,
		appService: appService,
	}
}

func (s *SecretService) GetSecret(secretID int64) (string, error) {
	result, err := s.dbClient.Context().Secret.Query().
		Where(secret.ID(secretID)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return "", server.NewError(http.StatusNotFound, fmt.Sprintf("secret with id %d not found", secretID))
		}
		return "", err
	}

	return result.Key, nil
}

func (s *SecretService) CreateSecret(appID int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error) {
	appModel, err := s.appService.GetAppByID(appID)
	if err != nil {
		return nil, err
	}

	alreadyExist, err := s.dbClient.Context().Secret.
		Query().
		Where(secret.Key(strings.ToUpper(secretModel.Key))).
		Exist(context.Background())

	if err != nil {
		return nil, err
	}

	if alreadyExist {
		return nil, server.NewError(http.StatusConflict, fmt.Sprintf("Secret %s for app already exist. ", secretModel.Key))
	}

	result, err := s.dbClient.Context().Secret.Create().
		SetKey(secretModel.Key).
		SetValue(secretModel.Value).
		SetAppID(appModel.ID).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	appSecretModel := new(model.AppSecretModel)
	appSecretModel.OriginalKey = secretModel.Key
	appSecretModel.Key = result.Key
	appSecretModel.SnippetURL = fmt.Sprintf("/apps/%d/secrets/%d/snippets", appID, result.ID)

	return appSecretModel, nil
}
