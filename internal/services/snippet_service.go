package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/internal/model"

	"github.com/internal/config"
)

const (
	Secret       = "secrets"
	GoFile       = "go.snippet"
	GoLanguage   = "Golang"
	GoClass      = "language-golang"
	GoInstall    = "go get -u gitlab.com/ikp_go-secrets"
	NodeFile     = "node.snippet"
	NodeLanguage = "Node.js"
	NodeClass    = "language-typescript"
	NodeInstall  = "npm install ikp_node-secrets@latest"
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
	secretName, err := s.secretService.GetSecret(secretID)
	if err != nil {
		return nil, err
	}

	for i, secret := range s.secrets {
		replaced := strings.ReplaceAll(secret.Code, "$PETS_APPNAME_SECRETKEY", secretName)
		s.secrets[i].Code = replaced
	}

	return s.secrets, nil
}

func buildSecrets(secrets []model.SnippetModel) []model.SnippetModel {
	secrets = append(secrets, buildSnippet(GoLanguage, Secret, GoFile, GoClass, GoInstall))
	secrets = append(secrets, buildSnippet(NodeLanguage, Secret, NodeFile, NodeClass, NodeInstall))
	return secrets
}

func buildSnippet(language string, snippetType string, file string, class string, install string) model.SnippetModel {
	snippet := new(model.SnippetModel)
	snippet.Class = class
	snippet.Language = language
	snippet.Install = install
	path := fmt.Sprintf("%s/%s/%s", config.String("snippets.folder"), snippetType, file)
	snippet.Code = getFileContent(path)

	return *snippet
}

func getFileContent(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		// fallback for test
		content, err = os.ReadFile("../../" + path)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return string(content)
}
