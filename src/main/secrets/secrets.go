package secrets

import (
	"errors"
	"fmt"
	"os"
)

type ISecretStore interface {
	GetSecret(key string) *Secret
}

type Secret struct {
	Key   string
	Value string
	Err   error
}

type SecretStore struct {
}

func NewSecretStore() *SecretStore {
	return &SecretStore{}
}

func (s *SecretStore) GetSecret(key string) *Secret {
	secret := new(Secret)
	secret.Key = key
	secret.Value = os.Getenv(key)
	if secret.Value == "" {
		secret.Err = errors.New(fmt.Sprintf("missing secret: %s", key))
	}

	return secret
}
