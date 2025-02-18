-- name: GetUserFromEmail :one
SELECT * FROM "User"
WHERE email = $1 LIMIT 1;

-- name: GetUserFromUsername :one
SELECT * FROM "User"
WHERE username = $1 LIMIT 1;

-- name: GetLiveClassFromId :one
SELECT * FROM "LiveClass"
WHERE id = $1 LIMIT 1;

-- name: GetLiveClassFromEmail :many
SELECT * FROM "LiveClass"
WHERE email = $1;