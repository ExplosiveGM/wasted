-- name: CreateUser :one
INSERT INTO users(login, code, code_expires_at) VALUES($1, $2, $3)
RETURNING *;