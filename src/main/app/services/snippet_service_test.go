package services_test

import (
	"testing"

	"github.com/src/main/app/model"
	"github.com/src/main/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSecretService struct {
	mock.Mock
}

func (m *MockSecretService) CreateSecret(int64, *model.CreateAppSecretModel) (*model.AppSecretModel, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockSecretService) GetSecret(int64) (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}

func TestSnippetService_GetSecrets(t *testing.T) {
	secretService := new(MockSecretService)
	secretService.On("GetSecret").Return(GetSecret())

	service := services.NewSnippetService(secretService)

	actual, err := service.GetSecrets(1)

	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Len(t, actual, 2)
	assert.Equal(t, actual[0].Language, string(services.GoLanguage))
	assert.Equal(t, actual[1].Language, string(services.NodeLanguage))
}

func GetSecret() (string, error) {
	return "MYSECRETKEY", nil
}
