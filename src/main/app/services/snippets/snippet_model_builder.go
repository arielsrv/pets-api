package snippets

import (
	"fmt"

	"github.com/src/main/app/config"
	"github.com/src/main/app/infrastructure/files"
	"github.com/src/main/app/model"
)

type ISnippetModelBuilder interface {
	IsSecret() *ISnippetSecretModelBuilder
	Build() model.SnippetModel
}

type SnippetModelBuilder struct {
	model model.SnippetModel
}

type ISnippetSecretModelBuilder interface {
	ForGo() *SnippetModelBuilder
	ForNode() *SnippetModelBuilder
}

func New() *SnippetModelBuilder {
	return &SnippetModelBuilder{model: model.SnippetModel{}}
}

func (builder *SnippetModelBuilder) IsSecret() *SnippetModelBuilder {
	builder.model.SnippetType = string(Secret)
	return builder
}

func (builder *SnippetModelBuilder) ForGo() *SnippetModelBuilder {
	builder.model.Language = "Golang"
	builder.model.Install = "go get -u gitlab.com/iskaypet/ikp_go-secrets"
	builder.model.File = "go.snippet"
	builder.model.Class = "language-golang"

	return builder
}

func (builder *SnippetModelBuilder) ForNode() *SnippetModelBuilder {
	builder.model.Language = "Node"
	builder.model.Install = "npm install ikp_node-secrets --save-dev"
	builder.model.File = "node.snippet"
	builder.model.Class = "language-typescript"

	return builder
}

func (builder *SnippetModelBuilder) Build() model.SnippetModel {
	path := fmt.Sprintf("%s/%s/%s",
		config.String("snippets.folder"),
		builder.model.SnippetType,
		builder.model.File)

	builder.model.Code = files.GetFileContent(path)

	return builder.model
}
