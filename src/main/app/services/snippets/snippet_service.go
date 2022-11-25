package snippets

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/src/main/app/services/secrets"

	"github.com/src/main/app/config"
	"github.com/src/main/app/model"
)

type SnippetType string

const (
	Secret SnippetType = "secrets"
)

type Language string

const (
	GoLanguage   Language = "Golang"
	NodeLanguage Language = "Node.js"
)

type File string

const (
	GoFile   File = "go.snippet"
	NodeFile File = "node.snippet"
)

type Class string

const (
	GoClass   Class = "language-golang"
	NodeClass Class = "language-typescript"
)

type Install string

const (
	GoInstall   Install = "go get -u gitlab.com/ikp_go-secrets" //
	NodeInstall Install = "npm install ikp_node-secrets@latest"
)

type ISnippetService interface {
	GetSecrets(secretID int64) ([]model.SnippetModel, error)
}

type SnippetService struct {
	secretService secrets.ISecretService
	snippets      map[string][]model.SnippetModel
}

type SnippetSecretModel struct {
	model.SnippetModel
}

type ISnippetStartupBuilder interface {
	BuildSecrets()
	Build() map[string][]model.SnippetModel
}

type SnippetStartupBuilder struct {
	snippets map[string][]model.SnippetModel
}

func (s SnippetStartupBuilder) BuildSecrets() {
	var secrets []model.SnippetModel
	secrets = append(secrets, buildSnippet(GoLanguage, Secret, GoFile, GoClass, GoInstall))
	secrets = append(secrets, buildSnippet(NodeLanguage, Secret, NodeFile, NodeClass, NodeInstall))
	s.snippets[string(Secret)] = secrets
}

func (s SnippetStartupBuilder) Build() map[string][]model.SnippetModel {
	s.snippets = make(map[string][]model.SnippetModel)
	s.BuildSecrets()

	return s.snippets
}

func NewSnippetService(secretService secrets.ISecretService) *SnippetService {
	snippetBuilder := new(SnippetStartupBuilder)
	snippets := snippetBuilder.Build()

	return &SnippetService{
		secretService: secretService,
		snippets:      snippets,
	}
}

func (s SnippetService) GetSecrets(secretID int64) ([]model.SnippetModel, error) {
	secretName, err := s.secretService.GetSecret(secretID)
	if err != nil {
		return nil, err
	}

	for i, secret := range s.snippets[string(Secret)] {
		replaced := strings.ReplaceAll(secret.Code, "$PETS_APPNAME_SECRETKEY", secretName)
		s.snippets[string(Secret)][i].Code = replaced
	}

	return s.snippets[string(Secret)], nil
}

func buildSnippet(language Language, snippetType SnippetType, file File, class Class, install Install) model.SnippetModel {
	snippet := new(model.SnippetModel)
	snippet.Class = string(class)
	snippet.Language = string(language)
	snippet.Install = string(install)
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
