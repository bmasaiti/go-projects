package storage

import (
	"fmt"
	"sync"

	"github.com/bmasaiti/go-projects/my-vault/internal/domain"
)

//template of my localdb ... How it should be according to me
type InMemorySecretRepo struct {
	mu sync.RWMutex
	secrets map[string]domain.Secret
}

//this func now creates an instance of the db
func NewInMemorySecretRepo() *InMemorySecretRepo {
    return &InMemorySecretRepo{
        secrets:  make(map[string]domain.Secret),
    }
}

//Now we can manipulate the db similar to how I was doing it with just a basic map.
// so we construct the secret object in the handler and we pass it to this function to actually add to our repository
// returns nothing if all goes well, and also now I have a localised mutex instead of global one (funny if I do repo.mu.Unlock its giving an error)
func (repo *InMemorySecretRepo) PutNewSecret(secret domain.Secret) error{
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.secrets[secret.Id]=secret
	return nil
}

func (repo *InMemorySecretRepo) GetScretsById(secretId string) (domain.Secret,error){
	repo.mu.Lock()
	defer repo.mu.Unlock()
	secretEntry, ok := repo.secrets[secretId]
	if !ok{
		return domain.Secret{}, fmt.Errorf("secret with ID '%s' not found", secretId)
	}
	//should this be a pointer instead .... not needed , since this is js a read , we can do with a copy
	return secretEntry, nil

}

func (repo *InMemorySecretRepo)DeleteSecretByID(secretId string) (string,error){
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.secrets[secretId]; exists {
		delete(repo.secrets, secretId)
		return secretId,nil
}
	return "", fmt.Errorf("secret with ID '%s' not found", secretId)
}

func (repo *InMemorySecretRepo)ListAllSecrets() ([]domain.Secret,error){
	repo.mu.Lock()
	defer repo.mu.Unlock()
	var secrets []domain.Secret
	if len(repo.secrets) == 0 {
		return []domain.Secret{},fmt.Errorf("no secrets found in the secrets store")
	}else{
		
	}

	for _, v := range repo.secrets {
		secrets = append(secrets, v)
	}
	return secrets,nil
}

