-- name: GetUser :one
SELECT * FROM users
WHERE uuid = ? LIMIT 1;

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
