package services

import (
	"fmt"
	"github.com/internal/model"
	"log"
	"os"
	"strings"

	"github.com/internal/config"
)

const (
	GoFile   = "go.snippets"
	NodeFile = "node.snippets"
)

const (
	Secret = "secrets"
)

const (
	GoLanguage   = "Golang"
	NodeLanguage = "Node.js"
)

const (
	GoClass   = "language-golang"
	NodeClass = "language-typescript"
)

type ISnippetService interface {
	GetSecrets(secretID int64) ([]model.SnippetModel, error)
}

type SnippetService struct {
	secretService ISecretService
	secrets       []model.SnippetModel
}

func NewSnippetService(secretService ISecretService) *SnippetService {
	var secrets []model.SnippetModel
	secrets = buildSecrets(secrets)

	return &SnippetService{
		secrets:       secrets,
		secretService: secretService,
	}
}

func (s SnippetService) GetSecrets(secretID int64) ([]model.SnippetModel, error) {
	secretName, appName, err := s.secretService.GetSecret(secretID)
	if err != nil {
		return nil, err
	}

	for i, secret := range s.secrets {
		snippetName := fmt.Sprintf("PETS_%s_%s", strings.ToUpper(appName), strings.ToUpper(secretName))
		replaced := strings.ReplaceAll(secret.Code, "$PETS_APPNAME_SECRETKEY", snippetName)
		s.secrets[i].Code = replaced
	}

	return s.secrets, nil
}

func buildSecrets(secrets []model.SnippetModel) []model.SnippetModel {
	secrets = append(secrets, buildSnippet(GoLanguage, Secret, GoFile, GoClass))
	secrets = append(secrets, buildSnippet(NodeLanguage, Secret, NodeFile, NodeClass))
	return secrets
}

func buildSnippet(language string, snippetType string, file string, class string) model.SnippetModel {
	snippet := new(model.SnippetModel)
	snippet.Class = class
	snippet.Language = language
	path := fmt.Sprintf("%s/%s/%s", config.String("snippets.folder"), snippetType, file)
	snippet.Code = getFileContent(path)

	return *snippet
}

func getFileContent(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	return string(content)
}
