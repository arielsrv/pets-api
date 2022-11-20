package secrets

type ISecretStore interface {
	GetSecret(key string) *Secret
}

type Secret struct {
	Key   string
	Value string
	Err   error
}
