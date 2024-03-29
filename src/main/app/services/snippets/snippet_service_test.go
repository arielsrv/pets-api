package snippets_test

import (
	"errors"
	"testing"

	"github.com/arielsrv/pets-api/src/main/app/services/snippets"

	"github.com/arielsrv/pets-api/src/main/app/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSecretService struct {
	mock.Mock
}

func (m *MockSecretService) CreateSecret(int64, *model.CreateSecretRequest) (*model.CreateSecretResponse, error) {
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

	service := snippets.NewSnippetService(secretService)

	actual, err := service.GetSecrets(1)

	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Len(t, actual, 2)
	assert.Equal(t, actual[0].Language, "Golang")
	assert.Equal(t, actual[1].Language, "Node")
}

func TestSnippetService_GetSecrets_NotFound(t *testing.T) {
	secretService := new(MockSecretService)
	secretService.On("GetSecret").Return(NotFound())

	service := snippets.NewSnippetService(secretService)

	actual, err := service.GetSecrets(1)

	assert.Error(t, err)
	assert.Nil(t, actual)
	assert.Equal(t, "secret with id 1 not found", err.Error())
}

func GetSecret() (string, error) {
	return "MYSECRETKEY", nil
}

func NotFound() (string, error) {
	return "MYSECRETKEY", errors.New("secret with id 1 not found")
}
