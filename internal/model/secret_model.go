package model

type SecretModel struct {
	Key         string `json:"key,omitempty"`
	RelativeUrl string `json:"url,omitempty"`
}
