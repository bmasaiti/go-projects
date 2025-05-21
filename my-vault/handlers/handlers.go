package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type Secret struct {
	Id    string
	Name  string
	KVMap map[string]string
}

var mu sync.Mutex
var db = make(map[string]Secret)

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
	Message string `json:"message"`
	Secrets []Secret `json:"Secrets"`

}
func GenerateUUID() string {
	newUUID, err := uuid.NewV7()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return err.Error()
	}
	return newUUID.String()
}

func NewSecretResponseObject(s Secret) GetSecretResponse {
	return GetSecretResponse{
		Message: "Retrieved secret object",
		Id:      s.Id,
		Name:    s.Name,
		KVMap:   s.KVMap,
	}
}

func NewCreateSecret(s CreateSecretRequest) Secret {
	return Secret{
		Id:    GenerateUUID(),
		Name:  s.Name,
		KVMap: s.KVMap,
	}
}

func BuildListSecretsResponse(s []Secret ) ListSecretsResponse{
	return ListSecretsResponse{
		Message: "Secrets fetched successfully",
		Secrets: s,
	}
}


func HandlePostSecret(w http.ResponseWriter, r *http.Request) {

	var secret Secret
	var secretRequestObject CreateSecretRequest

	err := json.NewDecoder(r.Body).Decode(&secretRequestObject)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	temp := NewCreateSecret(secretRequestObject)
	mu.Lock()
	db[temp.Id] = temp
	mu.Unlock()
	
	res := CreateSecretResponse{
		Id:      temp.Id,
		Name:    temp.Name,
		Message: fmt.Sprintf("Successfully created secret with secretId: %s and name: %s", temp.Id, temp.Name),
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(res)

	if err != nil {
		
		err := errors.New("unexpected internal error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Secret saved----------------------------------------- %s", secret)
}

func HandleGetSecretById(w http.ResponseWriter, r *http.Request) {
	//curl  -X POST -H "Content-Type: application/json" http://localhost:9000/secrets/234
	secretID := r.PathValue("secret_id")
	mu.Lock()
	secretEntry, ok := db[secretID]
	mu.Unlock()
	
	if !ok {
		http.Error(w, fmt.Sprintf("Secret with ID %s not found", secretID), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	getSecretResponse := NewSecretResponseObject(secretEntry)
	json.NewEncoder(w).Encode(getSecretResponse)
	fmt.Printf("Retrieved secret----------------------------------------- %s", getSecretResponse)
}

func HandleDeleteSecretById(w http.ResponseWriter, r *http.Request) {
	//curl  -X POST -H "Content-Type: application/json" http://localhost:9000/secrets/234
	secret_id := r.PathValue("secret_id")
	if _, exists := db[secret_id]; exists {
		mu.Lock()
		delete(db, secret_id)
		mu.Unlock()
		fmt.Println("Deleted secret with id -----------------------------------------", secret_id)
		fmt.Fprintf(w, "Secret with ID %s deleted successfully", secret_id)
	} else {
		http.Error(w, fmt.Sprintf("Secret with ID %s not found", secret_id), http.StatusNotFound)
		return
	}

}

func HandleListSecrets(w http.ResponseWriter, r *http.Request) {
	var secrets []Secret
	mu.Lock()
	if len(db) == 0 {
		http.Error(w, "No secrets found in the secrets store", http.StatusNotFound)
		return
	}
	
	for _, v := range db {
		secrets = append(secrets, v)
	}
	mu.Unlock()
	response := BuildListSecretsResponse(secrets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
