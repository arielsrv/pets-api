package snippets_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src/main/app/handlers/snippets"
	"github.com/src/main/app/model"
	"github.com/src/main/app/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSnippetService struct {
	mock.Mock
}

func (m *MockSnippetService) GetSecrets(int64) ([]model.Snippet, error) {
	args := m.Called()
	return args.Get(0).([]model.Snippet), args.Error(1)
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

func GetSecrets() ([]model.Snippet, error) {
	var snippetModel []model.Snippet

	var snippet1 model.Snippet
	snippet1.Language = "Golang"
	snippet1.Class = "language-golang"
	snippet1.Code = "main()"

	snippetModel = append(snippetModel, snippet1)

	var snippet2 model.Snippet
	snippet2.Language = "Node.js"
	snippet2.Class = "language-typescript"
	snippet2.Code = "console.log()"

	snippetModel = append(snippetModel, snippet2)

	return snippetModel, nil
}
