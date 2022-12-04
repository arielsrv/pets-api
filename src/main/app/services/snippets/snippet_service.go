package snippets

import (
	"strings"

	"github.com/src/main/app/services/secrets"

	"github.com/src/main/app/model"
)

type SnippetType string

const (
	Secret SnippetType = "secrets"
)

type ISnippetService interface {
	GetSecrets(secretID int64) ([]model.SnippetViewModel, error)
}

type SnippetService struct {
	secretService secrets.ISecretService
	snippets      map[string][]model.SnippetViewModel
}

type ISnippetStartupBuilder interface {
	BuildSecrets()
	Build() map[string][]model.SnippetViewModel
}

type SnippetStartupBuilder struct {
	snippets map[string][]model.SnippetViewModel
}

func (s SnippetStartupBuilder) BuildSecrets() {
	s.addOrUpdate(Secret, New().IsSecret().ForGo().Build())
	s.addOrUpdate(Secret, New().IsSecret().ForNode().Build())
}

func (s SnippetStartupBuilder) addOrUpdate(snippetType SnippetType, snippetModel model.SnippetViewModel) {
	s.snippets[string(snippetType)] = append(s.snippets[string(snippetType)], snippetModel)
}

func (s SnippetStartupBuilder) Build() map[string][]model.SnippetViewModel {
	s.snippets = make(map[string][]model.SnippetViewModel)
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

func (s SnippetService) GetSecrets(secretID int64) ([]model.SnippetViewModel, error) {
	secretName, err := s.secretService.GetSecret(secretID)
	if err != nil {
		return nil, err
	}

	var snippets []model.SnippetViewModel
	for _, secret := range s.snippets[string(Secret)] {
		snippet := new(model.SnippetViewModel)

		snippet.SnippetType = secret.SnippetType
		snippet.Language = secret.Language
		snippet.Install = secret.Install
		snippet.Class = secret.Class
		snippet.Code = strings.ReplaceAll(secret.Code, "$PETS_APPNAME_SECRETKEY", secretName)

		snippets = append(snippets, *snippet)
	}

	return snippets, nil
}
