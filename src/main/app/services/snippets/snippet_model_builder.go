package snippets

import (
	"fmt"

	"github.com/src/main/app/config"
	"github.com/src/main/app/infrastructure/files"
	"github.com/src/main/app/model"
)

type ISnippetModelBuilder interface {
	IsSecret() *ISnippetSecretModelBuilder
	Build() model.Snippet
}

type SnippetModelBuilder struct {
	model model.Snippet
}

type ISnippetSecretModelBuilder interface {
	ForGo() *SnippetModelBuilder
	ForNode() *SnippetModelBuilder
}

func New() *SnippetModelBuilder {
	return &SnippetModelBuilder{model: model.Snippet{}}
}

func (smb *SnippetModelBuilder) IsSecret() *SnippetModelBuilder {
	smb.model.SnippetType = string(Secret)
	return smb
}

func (smb *SnippetModelBuilder) ForGo() *SnippetModelBuilder {
	smb.model.Language = "Golang"
	smb.model.Install = "go get -u gitlab.com/iskaypet/ikp_go-secrets"
	smb.model.File = "go.snippet"
	smb.model.Class = "language-golang"

	return smb
}

func (smb *SnippetModelBuilder) ForNode() *SnippetModelBuilder {
	smb.model.Language = "Node"
	smb.model.Install = "npm install ikp_node-secrets --save-dev"
	smb.model.File = "node.snippet"
	smb.model.Class = "language-typescript"

	return smb
}

func (smb *SnippetModelBuilder) Build() model.Snippet {
	path := fmt.Sprintf("%s/%s/%s",
		config.String("snippets.folder"),
		smb.model.SnippetType,
		smb.model.File)

	smb.model.Code = files.GetFileContent(path)

	return smb.model
}
