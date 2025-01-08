// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

const getUserFromUsername = `-- name: GetUserFromUsername :one
SELECT id, username, password, name, email, picture, "authType" FROM "User"
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserFromUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserFromUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Name,
		&i.Email,
		&i.Picture,
		&i.AuthType,
	)
	return i, err
}
