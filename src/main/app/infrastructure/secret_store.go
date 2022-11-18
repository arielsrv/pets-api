package infrastructure

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

func (s SecretStore) GetSecret(string) *Secret {
	//TODO: @apineiro implement me
	panic("implement me")
}
