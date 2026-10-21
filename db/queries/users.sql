-- name: CreateUser :exec
INSERT INTO users (id, fullname, email, role, password_hash) VALUES($1, $2, $3, $4, $5);

-- name: FindByEmail :one
SELECT id, fullname, email, role, password_hash FROM users WHERE email = $1;