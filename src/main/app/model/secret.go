package model

type CreateSecretRequest struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type CreateSecretResponse struct {
	Key         string `json:"key,omitempty"`
	SnippetURL  string `json:"snippet_url,omitempty"`
	OriginalKey string `json:"original_key,omitempty"`
}
