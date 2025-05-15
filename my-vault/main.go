package main

import (
	"fmt"
	"net/http"
	"os"
	"my-vault/handlers"
	
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

func main() {

	fmt.Println("Starting Secrets Server--------------------------------------")
	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	router.HandleFunc("POST /v1/secrets", handlers.HandlePostSecret)
	router.HandleFunc("GET /v1/secrets", handlers.HandleListSecrets)
	router.HandleFunc("GET /v1/secrets/{secret_id}", handlers.HandleGetSecretById)
	router.HandleFunc("DELETE /v1/secrets/{secret_id}", handlers.HandleDeleteSecretById)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
        fmt.Println("Failed to listen and serve:", err)
        os.Exit(1)
    }
}

