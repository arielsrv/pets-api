package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/internal/config"
)

const (
	GoFile   = "go.snippet"
	NodeFile = "node.snippet"
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

type Snippet struct {
	Language string
	Class    string
	Code     string
}

type ISnippetService interface {
	GetSecrets(secretID int64) ([]Snippet, error)
}

type SnippetService struct {
	secretService ISecretService
	secrets       []Snippet
}

func NewSnippetService(secretService ISecretService) *SnippetService {
	var secrets []Snippet
	secrets = buildSecrets(secrets)

	return &SnippetService{
		secrets:       secrets,
		secretService: secretService,
	}
}

func (s SnippetService) GetSecrets(secretID int64) ([]Snippet, error) {
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

func buildSecrets(secrets []Snippet) []Snippet {
	secrets = append(secrets, buildSnippet(GoLanguage, Secret, GoFile, GoClass))
	secrets = append(secrets, buildSnippet(NodeLanguage, Secret, NodeFile, NodeClass))
	return secrets
}

func buildSnippet(language string, snippetType string, file string, class string) Snippet {
	snippet := new(Snippet)
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
