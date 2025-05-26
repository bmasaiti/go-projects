package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bmasaiti/go-projects/my-vault/internal/domain"
	"github.com/google/uuid"
)

//var db = storage.NewInMemorySecretRepo()
//return domain.Secret{}, fmt.Errorf("secret with ID '%s' not found", secretId)
//fmt.Errorf("no secrets found in the secrets store")
var ErrNotFound = errors.New("not found")
type SecretsRepository interface {
	PutNewSecret(secret domain.Secret) error
	GetScretsById(Id string) (domain.Secret, error)
	DeleteSecretByID(Id string) (string, error)
	ListAllSecrets() ([]domain.Secret, error)
}

// Couldn't mock the tests , need to inject the db
type SecretHandler struct {
	DB SecretsRepository
}

type CreateSecretRequest struct {
	Name  string            `json:"name"`
	KVMap map[string]string `json:"kv_map"`
}

type CreateSecretResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
	Name    string `json:"name"`
}
type GetSecretResponse struct {
	Message string            `json:"message"`
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	KVMap   map[string]string `json:"kv_map"`
}

type ListSecretsResponse struct {
	Message string          `json:"message"`
	Secrets []domain.Secret `json:"Secrets"`
}

func GenerateUUID() string {
	newUUID, err := uuid.NewV7()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return err.Error()
	}
	return newUUID.String()
}

func NewSecretResponseObject(s domain.Secret) GetSecretResponse {
	return GetSecretResponse{
		Message: "Retrieved secret object",
		Id:      s.Id,
		Name:    s.Name,
		KVMap:   s.KVMap,
	}
}

func NewCreateSecret(s CreateSecretRequest) domain.Secret {
	return domain.Secret{
		Id:    GenerateUUID(),
		Name:  s.Name,
		KVMap: s.KVMap,
	}
}

func BuildListSecretsResponse(s []domain.Secret) ListSecretsResponse {
	return ListSecretsResponse{
		Message: "Secrets fetched successfully",
		Secrets: s,
	}
}

func (h *SecretHandler) HandlePostSecret(w http.ResponseWriter, r *http.Request) {

	//var secret Secret
	var secretRequestObject CreateSecretRequest

	err := json.NewDecoder(r.Body).Decode(&secretRequestObject)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("ERROR: Unexpected behaviour: %v", err)
		return
	}
	// build new secretObject
	temp := NewCreateSecret(secretRequestObject)
	//err = db.PutNewSecret(temp)
	err = h.DB.PutNewSecret(temp)

	if err != nil {

		err := errors.New("unexpected internal error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR: Unexpted internal error: %v", err)
		return
	}

	res := CreateSecretResponse{
		Id:      temp.Id,
		Name:    temp.Name,
		Message: fmt.Sprintf("Successfully created secret with secretId: %s and name: %s", temp.Id, temp.Name),
	}
	
	log.Printf("INFO: Secret saved----------------------------------------- %s", temp.Id)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(res)
	if err != nil {

		err := errors.New("unexpected internal error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR: Unexpted internal error: %v", err)
		return
	}
}

func (h *SecretHandler) HandleGetSecretById(w http.ResponseWriter, r *http.Request) {
	//curl  -X POST -H "Content-Type: application/json" http://localhost:9000/secrets/234
	secretID := r.PathValue("secret_id")
	secretEntry, err := h.DB.GetScretsById(secretID)
	
	if err.Error() == ErrNotFound.Error(){
		http.Error(w, fmt.Sprintf("Secret with ID %s not found", secretID), http.StatusNotFound)
		return
	} else if err != nil{
		http.Error(w, "Unexpected internal error", http.StatusInternalServerError)
		log.Printf("ERROR: Unexpted internal error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	getSecretResponse := NewSecretResponseObject(secretEntry)
	json.NewEncoder(w).Encode(getSecretResponse)
	fmt.Printf("Retrieved secret----------------------------------------- %s", getSecretResponse)
}

func (h *SecretHandler) HandleDeleteSecretById(w http.ResponseWriter, r *http.Request) {
	//curl  -X POST -H "Content-Type: application/json" http://localhost:9000/secrets/234
	secretId := r.PathValue("secret_id")
	secret, err := h.DB.DeleteSecretByID(secretId)
	if err.Error()==ErrNotFound.Error(){
		http.Error(w, fmt.Sprintf("Secret with ID %s not found", secret), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to delete secret: %v", err)
		return
	}
	log.Println("Deleted secret with id -----------------------------------------", secret)
	fmt.Fprintf(w, "Secret with ID %s deleted successfully", secret)
}

func (h *SecretHandler) HandleListSecrets(w http.ResponseWriter, r *http.Request) {

	secrets, err := h.DB.ListAllSecrets()
	if err.Error() == ErrNotFound.Error(){
		http.Error(w,"No secrets found in the vault", http.StatusNotFound)
		return
	} else if err != nil{
		http.Error(w, "Unexpected internal server error", http.StatusInternalServerError)
		log.Printf("ERROR: Unexpted internal error: %v", err)
		return
	}
	response := BuildListSecretsResponse(secrets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
