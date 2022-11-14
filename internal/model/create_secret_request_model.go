package model

type CreateSecretRequestModel struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
