package domain
type Secret struct {
	Id    string
	Name  string
	KVMap map[string]string
}

type SecretsRepository interface{
	PutNewSecret(secret Secret) (Secret,error)
	GetScretsById(Id string)(Secret,error)
	DeleteSecretByID(Id string)(string,error)
	ListAllSecrets()([]Secret,error)
}

