package model

type SnippetModel struct {
	Language string `json:"language,omitempty"`
	Class    string `json:"class,omitempty"`
	Code     string `json:"code,omitempty"`
}
