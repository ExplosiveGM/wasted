-- name: CheckUserByRefreshCode :one
SELECT * FROM users
WHERE refresh_token = $1 AND refresh_token_expires_at >= NOW() LIMIT 1;
