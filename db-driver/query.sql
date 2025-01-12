-- name: GetUserFromUsername :one
SELECT * FROM "User"
WHERE username = $1 LIMIT 1;

-- name: GetUserFromEmail :one
SELECT * FROM "User"
WHERE email = $1 LIMIT 1;