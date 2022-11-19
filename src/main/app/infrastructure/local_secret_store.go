package infrastructure

import (
	"fmt"
	"github.com/src/main/app/config"
)

type LocalSecretStore struct {
}

func NewLocalSecretStore() *LocalSecretStore {
	return &LocalSecretStore{}
}

func (l *LocalSecretStore) GetSecret(key string) *Secret {
	secret := new(Secret)
	secret.Key = key
	secret.Value = config.String(key)
	if secret.Value == "" {
		secret.Err = fmt.Errorf("missing secret: %s", key)
	}

	return secret
}
