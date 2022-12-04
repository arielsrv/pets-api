package model

type Secret struct {
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	SnippetURL  string `json:"snippet_url,omitempty"`
	OriginalKey string `json:"original_key,omitempty"`
}
