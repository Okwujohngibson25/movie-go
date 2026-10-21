-- name: CreateScreen :exec
INSERT INTO screens (id, title, capacity) VALUES ($1, $2, $3);

-- name: GetAllScreens :many
SELECT id, title, capacity FROM screens;

-- name: GetScreensByID :one
SELECT id, title, capacity FROM screens WHERE id = $1;

-- name: UpdateScreensByID :execrows
UPDATE screens SET title = $1 WHERE id = $2;