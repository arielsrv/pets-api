package model

type AppSecretModel struct {
	Key      string         `json:"key,omitempty"`
	Snippets []SnippetModel `json:"snippets,omitempty"`
}

type SnippetModel struct {
	Language string `json:"language,omitempty"`
	CodeUrl  string `json:"code_url,omitempty"`
}
