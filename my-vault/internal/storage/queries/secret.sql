
-- name: PutNewSecret :exec
INSERT INTO secrets (id, name, kv_data) VALUES ($1, $2, $3);

-- name: GetScretsById :one
SELECT id, name, kv_data FROM secrets WHERE id = $1;

-- name: DeleteSecretByID :exec
DELETE FROM secrets WHERE id = $1;

-- name: ListAllSecrets :many
SELECT id, name, kv_data FROM secrets;