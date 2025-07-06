package storage

import (
	"database/sql"
	"fmt"
	"log"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgressDBConnection(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	if err := CreateTables(db); err != nil {
		log.Fatalf("Could not create tables: %v", err)
	}

	return db, nil
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
