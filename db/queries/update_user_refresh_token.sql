-- name: UpdateUserRefreshToken :exec
UPDATE users
SET refresh_token = $1,refresh_token_expires_at = $2
WHERE id = $3;