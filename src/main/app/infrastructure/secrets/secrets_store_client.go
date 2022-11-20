package secrets

type SecretStore struct {
}

func NewSecretStore() *SecretStore {
	return &SecretStore{}
}

func (s SecretStore) GetSecret(string) *Secret {
	//TODO: @apineiro implement me
	panic("implement me")
}
