// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (
  uuid, name, email, password, create_time
) VALUES (
  ?, ?, ?, ?, ?
)
`

type CreateUserParams struct {
	Uuid       string       `json:"uuid"`
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Password   string       `json:"password"`
	CreateTime sql.NullTime `json:"create_time"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser,
		arg.Uuid,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.CreateTime,
	)
}

const getUser = `-- name: GetUser :one
SELECT uuid, name, email, password, create_time FROM users
WHERE uuid = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, uuid string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, uuid)
	var i User
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreateTime,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
  SET password = ?
WHERE uuid = ?
`

type UpdateUserParams struct {
	Password string `json:"password"`
	Uuid     string `json:"uuid"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser, arg.Password, arg.Uuid)
	return err
}
