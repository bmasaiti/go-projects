package main

import (
	"database/sql"
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


type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	}


func LoadDatabaseConfig() DatabaseConfig {
		return DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "adminYourPass"),
			DBName:   getEnv("DB_NAME", "secrets_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		}
	}

func getEnv(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
func NewPostgressDBConnection(config DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db,nil
}


func CreateTables(db *sql.DB) error {

	createSecretsTable := `
	CREATE TABLE IF NOT EXISTS secrets (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		kv_data JSONB NOT NULL
	);
	`
	if _, err := db.Exec(createSecretsTable); err != nil {
		return fmt.Errorf("failed to create secrets table: %w", err)
	}
	log.Println("Database tables created successfully")
	return nil
}

func main() {
	
	
	

	// repo := storage.NewInMemorySecretRepo()

    // secretHandler := &handlers.SecretHandler{
    //     DB: repo,
    // }



	config := LoadDatabaseConfig()
    db, err := NewPostgressDBConnection(config)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    CreateTables(db)
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

	router.HandleFunc("POST /v1/secrets", func(w http.ResponseWriter, r *http.Request){
		secretHandler.HandlePostSecret(w, r)
	})
	router.HandleFunc("GET /v1/secrets/{secret_id}", func(w http.ResponseWriter, r *http.Request){
		secretHandler.HandleGetSecretById(w, r)
	})
	router.HandleFunc("GET /v1/secrets", func(w http.ResponseWriter, r *http.Request){
		secretHandler.HandleListSecrets(w, r)
	})
	router.HandleFunc("DELETE /v1/secrets/{secret_id}", func(w http.ResponseWriter, r *http.Request){
		secretHandler.HandleDeleteSecretById(w, r)
	})
	
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
        fmt.Println("Failed to listen and serve:", err)
        os.Exit(1)
    }

}

