package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/internal/server"

	"github.com/internal/model"

	"github.com/internal/handlers"
	"github.com/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSnippetService struct {
	mock.Mock
}

func (m *MockSnippetService) GetSecrets(int64) ([]model.SnippetModel, error) {
	args := m.Called()
	return args.Get(0).([]model.SnippetModel), args.Error(1)
}

func TestSnippetHandler_GetSnippet(t *testing.T) {
	snippetService := new(MockSnippetService)
	snippetService.On("GetSecrets").Return(GetSecrets())
	snippetHandler := handlers.NewSnippetHandler(snippetService)

	app := server.New()

	app.Add(http.MethodGet, "/apps/:appId/secrets/:secretId/snippets", snippetHandler.GetSnippet)

	request := httptest.NewRequest(http.MethodGet, "/apps/1/secrets/2/snippets", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)

	assert.NotEmpty(t, string(body))
}

func GetSecrets() ([]model.SnippetModel, error) {
	var snippets []model.SnippetModel

	var snippet1 model.SnippetModel
	snippet1.Language = services.GoLanguage
	snippet1.Class = services.GoClass
	snippet1.Code = "main()"

	snippets = append(snippets, snippet1)

	var snippet2 model.SnippetModel
	snippet2.Language = services.NodeLanguage
	snippet2.Class = services.NodeClass
	snippet2.Code = "console.log()"

	snippets = append(snippets, snippet2)

	return snippets, nil
}
