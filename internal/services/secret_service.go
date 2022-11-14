package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ent"
	"github.com/ent/secret"
	"github.com/internal/clients/gitlab"
	"github.com/internal/infrastructure"
	"github.com/internal/model"
	"github.com/internal/shared"
)

type ISecretService interface {
	SaveSecret(appId int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error)
	GetSecret(secretID int64) (string, string, error)
}

type SecretService struct {
	client     gitlab.IGitLabClient
	dataAccess *infrastructure.DataAccessService
}

func NewSecretService(client gitlab.IGitLabClient, dataAccess *infrastructure.DataAccessService) *SecretService {
	return &SecretService{client: client, dataAccess: dataAccess}
}

func (s *SecretService) GetSecret(secretID int64) (string, string, error) {
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

func (s *SecretService) SaveSecret(appId int64, secretModel *model.CreateAppSecretModel) (*model.AppSecretModel, error) {
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
