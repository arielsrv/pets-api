package snippets_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src/main/app/handlers/snippets"
	snippets2 "github.com/src/main/app/services/snippets"

	"github.com/src/main/app/model"
	"github.com/src/main/app/server"
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
	snippetHandler := snippets.NewSnippetHandler(snippetService)

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
	snippet1.Language = string(snippets2.GoLanguage)
	snippet1.Class = string(snippets2.GoClass)
	snippet1.Code = "main()"

	snippets = append(snippets, snippet1)

	var snippet2 model.SnippetModel
	snippet2.Language = string(snippets2.NodeLanguage)
	snippet2.Class = string(snippets2.NodeClass)
	snippet2.Code = "console.log()"

	snippets = append(snippets, snippet2)

	return snippets, nil
}
