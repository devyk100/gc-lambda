-- name: GetUserFromEmail :one
SELECT * FROM "User"
WHERE email = $1 LIMIT 1;