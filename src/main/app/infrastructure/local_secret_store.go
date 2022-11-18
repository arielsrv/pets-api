package infrastructure

import (
	"fmt"
	"os"
)

type LocalSecretStore struct {
}

func NewLocalSecretStore() *LocalSecretStore {
	return &LocalSecretStore{}
}

func (l *LocalSecretStore) GetSecret(key string) *Secret {
	secret := new(Secret)
	secret.Key = key
	secret.Value = os.Getenv(key)
	if secret.Value == "" {
		secret.Err = fmt.Errorf("missing secret: %s", key)
	}

	return secret
}
