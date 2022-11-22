package model

type AppSecretModel struct {
	Key         string `json:"key,omitempty"`
	SnippetURL  string `json:"snippet_url,omitempty"`
	OriginalKey string `json:"original_key,omitempty"`
}
