package main

import (
	"encoding/json"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Create a HTTP server that allows you to create, delete, read and list “secrets”.
//  The following paths and verbs should be used for each ({} denotes a path variable):
// Create: POST /secrets
// Delete: DELETE /secrets/{secret_id}
// Read: GET /secrets/{secret_id}
// List: GET /secrets
// Evertything can be stored in-memory and no need for authentication for now.
// Each endpoint should accept and return a JSON formatted string (fields are up to you).

// curl \
//     --header "X-Vault-Token: ..." \
//     --request POST \
//     --data @payload.json \
//     https://127.0.0.1:8200/v1/secret/data/my-secret

//curl -d "@request_payload.json" -X POST  http://localhost:8000/secrets

//TODO: Create a db object that can be swapped out for a  real db 
//TODO: add decent logging.
//TODO: Proper error handling.
//TODO: encrypting secret object. 
type Secret struct {
	Id    int               `json:id`
	Name  string            `json:name`
	KVMap map[string]string `json:"kv_map"`
}

var db = make(map[int]Secret)
var db_ptr = &db

func main() {

	fmt.Println("Starting Secrets Server--------------------------------------")
	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	router.HandleFunc("POST /v1/secrets", HandlePostSecret)
	router.HandleFunc("GET /v1/secrets", HandleListSecrets)
	router.HandleFunc("GET /v1/secrets/{secret_id}", HandleGetSecretById)
	router.HandleFunc("DELETE /v1/secrets/{secret_id}", HandleDeleteSecretById)

	server.ListenAndServe()
}

func HandlePostSecret(w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(100000)

	var secret Secret
	err := json.NewDecoder(r.Body).Decode(&secret)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	secret.Id = id
	db[secret.Id] = secret

	fmt.Println("Secret saved----------------------------------------- %s", &secret)
}

func HandleGetSecretById(w http.ResponseWriter, r *http.Request) {
	//curl  -X POST -H "Content-Type: application/json" http://localhost:9000/secrets/234
	secret_id, err := strconv.Atoi(html.EscapeString(r.PathValue("secret_id")))

	if err != nil {
		fmt.Errorf(err.Error())
		fmt.Println("Failed to convert string to integer")
		return
	}
	secret_entry := db[secret_id]
	encoder := json.NewEncoder(w)
	encoded_response := encoder.Encode(secret_entry)
	data, err := json.Marshal(encoded_response)
	if err!= nil {
		fmt.Errorf(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))

	fmt.Printf("Retrieved secret saved----------------------------------------- %s", encoded_response)
}

func HandleDeleteSecretById(w http.ResponseWriter, r *http.Request) {
	//curl  -X POST -H "Content-Type: application/json" http://localhost:9000/secrets/234
	secret_id, err := strconv.Atoi(html.EscapeString(r.PathValue("secret_id")))
	if err != nil {
		fmt.Println("Failed to convert string to integer")
		return
	}

	delete(db, secret_id)

	fmt.Printf("Deleted secret with id -----------------------------------------", secret_id)
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
