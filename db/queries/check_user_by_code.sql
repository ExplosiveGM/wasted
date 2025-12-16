-- name: CheckUserByCode :one
SELECT * FROM users
WHERE login = $1 AND code = $2 AND code_expires_at >= NOW() LIMIT 1;
