-- name: UpdateUser :exec
UPDATE users
SET code = $1, code_expires_at = $2
WHERE id = $3;