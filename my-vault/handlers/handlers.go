package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"sync"
)

type Secret struct {
	Id    string 
	Name  string  
	KVMap map[string]string 
	
}
var mu sync.Mutex
var db = make(map[string]Secret)


type CreateSecretRequest  struct {
	Name  string    `json:"name"`
	KVMap map[string]string `json:"kv_map"`
	
}

type CreateSecretResponse  struct {
	Message string `json:"message"`
	Id string `json:"id"`
	Name  string    `json:"name"`
	
	
}
type GetSecretResponse struct {
	Message string `json:"message"`
	Id    string  `json:"id"`
	Name  string    `json:"name"`
	KVMap map[string]string `json:"kv_map"`

}

func GenerateUUID() string{
	newUUID, err := uuid.NewV7()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		panic(err)
	}
	return  newUUID.String()
}

func NewSecretResponseObject(s Secret) GetSecretResponse {
	return GetSecretResponse{
		Message: "Retrieved secret object",
		Id : s.Id,
		Name: s.Name,
		KVMap: s.KVMap,
	}
}

// func NewCreateSecret(s CreateSecretRequest ) Secret {
// 	return Secret{
// 		Id: GenerateUUID(),
// 		Name: s.Name,
// 		KVMap: s.KVMap,

// 	}
// }

func HandlePostSecret(w http.ResponseWriter, r *http.Request) {

	
	var secret Secret
	var secretRequestObject CreateSecretRequest

	err := json.NewDecoder(r.Body).Decode(&secretRequestObject)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secret.Id = GenerateUUID()
	secret.Name = secretRequestObject.Name
	secret.KVMap = secretRequestObject.KVMap


	mu.Lock()
	//db[secret.Id] =  NewCreateSecret(secretRequestObject)
	db[secret.Id] = secret
	mu.Unlock()
	defer r.Body.Close()

	 res:= CreateSecretResponse{
			Id: secret.Id,
			Name: secret.Name,
			Message: fmt.Sprintf("Successfully created secret with secretId: %s and name: %s", secret.Id, secret.Name),
	 }
	
	encoder := json.NewEncoder(w)
	encoded_response := encoder.Encode(res )
	data, err := json.Marshal(encoded_response)
	if err!= nil {
		fmt.Errorf("%s", err.Error())
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
	fmt.Println("Secret saved----------------------------------------- %s", secret)
}

func HandleGetSecretById(w http.ResponseWriter, r *http.Request) {
	//curl  -X POST -H "Content-Type: application/json" http://localhost:9000/secrets/234
	secret_id := r.PathValue("secret_id")

	// if secret_id != nil {
	// 	fmt.Errorf(err.Error())
	// 	fmt.Println("Failed to convert string to integer")
	// 	return
	// }
	secret_entry := db[secret_id]
	encoder := json.NewEncoder(w)

	getSecretResponse := NewSecretResponseObject(secret_entry)
	
	encoded_response := encoder.Encode(getSecretResponse )
	data, err := json.Marshal(encoded_response)
	if err!= nil {
		fmt.Errorf("%s", err.Error())
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))

	fmt.Printf("Retrieved secret saved----------------------------------------- %s", encoded_response)
}

func HandleDeleteSecretById(w http.ResponseWriter, r *http.Request) {
	//curl  -X POST -H "Content-Type: application/json" http://localhost:9000/secrets/234
	secret_id := r.PathValue("secret_id")
	// if err != nil {
	// 	fmt.Println("Failed to convert string to integer")
	// 	return
	// }
	mu.Lock()
	delete(db, secret_id)
	mu.Unlock()

	fmt.Println("Deleted secret with id -----------------------------------------", secret_id)
}

func HandleListSecrets(w http.ResponseWriter, r *http.Request) {
	var secrets []Secret
	for _, v := range db {
		secrets = append(secrets, v)
	}
	data, err := json.Marshal(secrets)
	if err!= nil {
		fmt.Errorf(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}


