-- name: GetUserFromUsername :one
SELECT * FROM "User"
WHERE username = $1 LIMIT 1;