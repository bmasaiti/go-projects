package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/bmasaiti/go-projects/my-vault/internal/domain"
)

type PostgresSecretRepo struct {
	DB *sql.DB
}

func NewPostgresSecretRepo(db *sql.DB) *PostgresSecretRepo {
	return &PostgresSecretRepo{DB: db}
}

func (repo *PostgresSecretRepo) PutNewSecret(secret domain.Secret) error {
	query := `
		INSERT INTO secrets (id, name, kv_data) 
		VALUES ($1, $2, $3)
	`
	kvDataJSON, err := json.Marshal(secret.KVMap)
	if err != nil {
		return fmt.Errorf("failed to marshal kv_data: %w", err)
	}

	_, err = repo.DB.Exec(query, secret.Id, secret.Name, kvDataJSON)
	if err != nil {
		return fmt.Errorf("failed to insert secret: %w", err)
	}

	return nil
}

func (repo *PostgresSecretRepo) GetScretsById(secretId string) (domain.Secret, error) {
	query := `
		SELECT id, name, kv_data
		FROM secrets 
		WHERE id = $1
	`
	var secret domain.Secret

	var kvDataJSON []byte

	err := repo.DB.QueryRow(query, secretId).Scan(
		&secret.Id,
		&secret.Name,
		&kvDataJSON,
	)
	if err == sql.ErrNoRows {
		return domain.Secret{}, fmt.Errorf("secret with id %s not found", secretId)
	}
	if err != nil {
		return domain.Secret{}, fmt.Errorf("failed to query secret: %w", err)
	}

	if err := json.Unmarshal(kvDataJSON, &secret.KVMap); err != nil {
		return domain.Secret{}, fmt.Errorf("failed to unmarshal kv_data: %w", err)
	}
	return secret, nil

}

func (repo *PostgresSecretRepo) DeleteSecretByID(secretId string) (string, error) {
	query := `
		DELETE FROM secrets 
		WHERE id = $1
	`
	var secret domain.Secret

	var kvDataJSON []byte

	err := repo.DB.QueryRow(query, secretId).Scan(
		&secret.Id,
		&secret.Name,
		&kvDataJSON,
	)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("secret with id %s not found", secretId)
	}
	if err != nil {
		return "", fmt.Errorf("failed to delete secret: %w", err)
	}
	return secretId, nil
}

func (repo *PostgresSecretRepo) ListAllSecrets() ([]domain.Secret, error) {
	query := `
		SELECT id, name, kv_data
		FROM secrets
	`
	var secrets []domain.Secret

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query secrets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var secret domain.Secret
		var kvDataJSON []byte

		err := rows.Scan(
			&secret.Id,
			&secret.Name,
			&kvDataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret row: %w", err)
		}

		if err := json.Unmarshal(kvDataJSON, &secret.KVMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal kv_data: %w", err)
		}

		secrets = append(secrets, secret)
	}
	return secrets, nil
}
