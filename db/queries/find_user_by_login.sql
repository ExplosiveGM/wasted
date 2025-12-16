-- name: FindUserByLogin :one
SELECT * FROM users
WHERE login = $1 LIMIT 1;