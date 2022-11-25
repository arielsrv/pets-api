package services_test

import (
	"context"
	"testing"

	"github.com/src/main/app/infrastructure/database"

	"github.com/src/main/app/model"
	"github.com/src/main/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) GetGroups() ([]model.AppGroupModel, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockAppService) CreateApp(*model.CreateAppModel) (*model.AppModel, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockAppService) GetAppTypes() ([]model.AppTypeModel, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockAppService) GetAppByName(string) (*model.AppModel, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockAppService) GetAppByID(int64) (*model.AppModel, error) {
	args := m.Called()
	return args.Get(0).(*model.AppModel), args.Error(1)
}

func TestSecretService_GetSecret(t *testing.T) {
	dbClient := database.NewDBClient(database.NewSQLiteClient(t))
	dbClient.Context()
	defer dbClient.Close()

	dbClient.AppType.Create().SetID(1).SetName("backend").Save(context.Background())
	dbClient.App.Create().SetName("customers-api").SetProjectId(1).SetAppTypeID(1).Save(context.Background())
	dbClient.Secret.Create().SetKey("PETS_CUSTOMERS-API_MYSECRETKEY").SetValue("MYSECRETVALUE").SetAppID(1).Save(context.Background())

	appService := new(MockAppService)

	service := services.NewSecretService(dbClient, appService)
	actual, err := service.GetSecret(1)

	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, "PETS_CUSTOMERS-API_MYSECRETKEY", actual)
}

func TestSecretService_GetSecret_NotFound(t *testing.T) {
	dbClient := database.NewDBClient(database.NewSQLiteClient(t))
	dbClient.Context()
	defer dbClient.Close()

	appService := new(MockAppService)

	service := services.NewSecretService(dbClient, appService)
	actual, err := service.GetSecret(2)

	assert.Error(t, err)
	assert.NotNil(t, actual)
}

func TestSecretService_CreateSecret(t *testing.T) {
	dbClient := database.NewDBClient(database.NewSQLiteClient(t))
	dbClient.Context()
	defer dbClient.Close()

	dbClient.AppType.Create().SetID(1).SetName("backend").Save(context.Background())
	dbClient.App.Create().SetName("customers-api").SetProjectId(1).SetAppTypeID(1).Save(context.Background())

	appService := new(MockAppService)
	appService.On("GetAppByID").Return(GetApp())

	service := services.NewSecretService(dbClient, appService)
	secretModel := new(model.CreateAppSecretModel)
	secretModel.Key = "MYSECRETKEY"
	secretModel.Value = "MYSECRETVALUE"
	actual, err := service.CreateSecret(1, secretModel)

	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, "MYSECRETKEY", actual.Key)
	assert.Equal(t, "/apps/1/secrets/1/snippets", actual.SnippetURL)
}

func TestSecretService_CreateSecret_Conflict(t *testing.T) {
	dbClient := database.NewDBClient(database.NewSQLiteClient(t))
	dbClient.Context()
	defer dbClient.Close()

	dbClient.AppType.Create().SetID(1).SetName("backend").Save(context.Background())
	dbClient.App.Create().SetName("customers-api").SetProjectId(1).SetAppTypeID(1).Save(context.Background())

	appService := new(MockAppService)
	appService.On("GetAppByID").Return(GetApp())

	service := services.NewSecretService(dbClient, appService)
	secretModel := new(model.CreateAppSecretModel)
	secretModel.Key = "MYSECRETKEY"
	secretModel.Value = "MYSECRETVALUE"
	actual, err := service.CreateSecret(1, secretModel)

	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, "MYSECRETKEY", actual.Key)
	assert.Equal(t, "/apps/1/secrets/1/snippets", actual.SnippetURL)

	conflict, err := service.CreateSecret(1, secretModel)

	assert.Error(t, err)
	assert.Nil(t, conflict)
}

func GetApp() (*model.AppModel, error) {
	appModel := new(model.AppModel)
	appModel.ID = 1
	appModel.Name = "MyApp"
	appModel.URL = "/apps/1/secrets/2/snippets"

	return appModel, nil
}
