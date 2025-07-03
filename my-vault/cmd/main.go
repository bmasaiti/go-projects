package main

import (
	
	"fmt"
	"log"
	"net/http"
	"os"
	 _ "github.com/lib/pq"
	"github.com/bmasaiti/go-projects/my-vault/internal/handlers"
	"github.com/bmasaiti/go-projects/my-vault/internal/storage"
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


func getEnv(key, defaultValue string) string {
	return cmp.Or(os.Getenv(key), defaultValue)
}

func main() {
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "admin"),
		getEnv("DB_PASSWORD", "adminYourPass"),
		getEnv("DB_NAME", "secrets_db"),
		getEnv("DB_SSLMODE", "disable"),
	)
	db, err := storage.NewPostgressDBConnection(connStr)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	
	defer db.Close()
	log.Println("Successfully connected to the database!")

	if err := storage.CreateTables(db); err != nil {
		log.Fatalf("Could not create tables: %v", err)
	}


	postgres:= storage.NewPostgresSecretRepo(db)

	secretHandler := &handlers.SecretHandler{
        DB: postgres,
    }


	fmt.Println("Starting Secrets Server--------------------------------------")
	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	router.HandleFunc("POST /v1/secrets", func(w http.ResponseWriter, r *http.Request){secretHandler.HandlePostSecret(w, r)})
	router.HandleFunc("GET /v1/secrets/{secret_id}", func(w http.ResponseWriter, r *http.Request){secretHandler.HandleGetSecretById(w, r)})
	router.HandleFunc("GET /v1/secrets", func(w http.ResponseWriter, r *http.Request){secretHandler.HandleListSecrets(w, r)})
	router.HandleFunc("DELETE /v1/secrets/{secret_id}", func(w http.ResponseWriter, r *http.Request){secretHandler.HandleDeleteSecretById(w, r)})
	
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
        fmt.Println("Failed to listen and serve:", err)
        os.Exit(1)
    }

}

