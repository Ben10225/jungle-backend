-- name: GetUserByUuid :one
SELECT * FROM users
WHERE uuid = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: GetUserByEmailAndPwd :one
SELECT * FROM users 
WHERE email = ? AND password = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (
  uuid, name, email, password, create_time
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: UpdateUser :exec
UPDATE users
  SET password = ?
WHERE uuid = ?;

