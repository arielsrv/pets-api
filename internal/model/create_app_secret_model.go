package model

type CreateAppSecretModel struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}