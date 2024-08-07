package secrets_test

import (
	"context"
	"testing"

	"github.com/arielsrv/pets-api/src/main/app/infrastructure/database"
	"github.com/arielsrv/pets-api/src/main/app/model"
	"github.com/arielsrv/pets-api/src/main/app/services/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) GetGroups() ([]model.AppGroupResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockAppService) CreateApp(*model.CreateAppRequest) (*model.CreateAppResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockAppService) GetAppTypes() ([]model.AppTypeResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockAppService) GetAppByName(string) (*model.AppResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockAppService) GetAppByID(int64) (*model.AppResponse, error) {
	args := m.Called()
	return args.Get(0).(*model.AppResponse), args.Error(1)
}

func TestSecretService_GetSecret(t *testing.T) {
	dbClient := database.NewDBClient(database.NewSQLiteClient(t))
	dbClient.Context()
	defer func(dbClient *database.DBClient) {
		err := dbClient.Close()
		require.NoError(t, err)
	}(dbClient)

	appType, err := dbClient.AppType.Create().
		SetID(1).SetName("backend").Save(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, appType)

	app, err := dbClient.App.Create().
		SetName("customers-api").SetExternalGitlabProjectID(1).SetAppTypeID(1).Save(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, app)

	secret, err := dbClient.Secret.Create().
		SetKey("PETS_CUSTOMERS-API_MYSECRETKEY").SetValue("MYSECRETVALUE").SetAppID(1).Save(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, secret)

	appService := new(MockAppService)

	service := secrets.NewSecretService(dbClient, appService)
	actual, err := service.GetSecret(1)

	require.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, "PETS_CUSTOMERS-API_MYSECRETKEY", actual)
}

func TestSecretService_GetSecret_NotFound(t *testing.T) {
	dbClient := database.NewDBClient(database.NewSQLiteClient(t))
	dbClient.Context()
	defer func(dbClient *database.DBClient) {
		err := dbClient.Close()
		require.NoError(t, err)
	}(dbClient)

	appService := new(MockAppService)

	service := secrets.NewSecretService(dbClient, appService)
	actual, err := service.GetSecret(2)

	require.Error(t, err)
	assert.NotNil(t, actual)
}

func TestSecretService_CreateSecret(t *testing.T) {
	dbClient := database.NewDBClient(database.NewSQLiteClient(t))
	dbClient.Context()
	defer func(dbClient *database.DBClient) {
		err := dbClient.Close()
		require.NoError(t, err)
	}(dbClient)

	appType, err := dbClient.AppType.Create().
		SetID(1).SetName("backend").Save(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, appType)

	app, err := dbClient.App.Create().
		SetName("customers-api").SetExternalGitlabProjectID(1).SetAppTypeID(1).Save(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, app)

	appService := new(MockAppService)
	appService.On("GetAppByID").Return(GetApp())

	service := secrets.NewSecretService(dbClient, appService)
	secretModel := new(model.CreateSecretRequest)
	secretModel.Key = "MYSECRETKEY"
	secretModel.Value = "MYSECRETVALUE"
	actual, err := service.CreateSecret(1, secretModel)

	require.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, "MYSECRETKEY", actual.Key)
	assert.Equal(t, "/apps/1/secrets/1/snippets", actual.SnippetURL)
}

func TestSecretService_CreateSecret_Conflict(t *testing.T) {
	dbClient := database.NewDBClient(database.NewSQLiteClient(t))
	dbClient.Context()
	defer func(dbClient *database.DBClient) {
		err := dbClient.Close()
		require.NoError(t, err)
	}(dbClient)

	appType, err := dbClient.AppType.Create().
		SetID(1).SetName("backend").Save(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, appType)

	app, err := dbClient.App.Create().
		SetName("customers-api").SetExternalGitlabProjectID(1).SetAppTypeID(1).Save(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, app)

	appService := new(MockAppService)
	appService.On("GetAppByID").Return(GetApp())

	service := secrets.NewSecretService(dbClient, appService)
	secretModel := new(model.CreateSecretRequest)
	secretModel.Key = "MYSECRETKEY"
	secretModel.Value = "MYSECRETVALUE"
	actual, err := service.CreateSecret(1, secretModel)

	require.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, "MYSECRETKEY", actual.Key)
	assert.Equal(t, "/apps/1/secrets/1/snippets", actual.SnippetURL)

	conflict, err := service.CreateSecret(1, secretModel)

	require.Error(t, err)
	assert.Nil(t, conflict)
}

func GetApp() (*model.AppResponse, error) {
	appModel := new(model.AppResponse)
	appModel.ID = 1
	appModel.Name = "MyApp"
	appModel.URL = "/apps/1/secrets/2/snippets"

	return appModel, nil
}
