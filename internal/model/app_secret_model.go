package model

type AppSecretModel struct {
	Key         string `json:"key,omitempty"`
	RelativeUrl string `json:"url,omitempty"`
}
