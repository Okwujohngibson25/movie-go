-- name: CreateMovie :one
INSERT INTO movies (id, title, description, duration_minutes, rating) VALUES($1, $2, $3, $4, $5) ON CONFLICT (title) DO NOTHING RETURNING *;

-- name: GetAllMovies :many
SELECT id, title, description, duration_minutes, rating FROM movies WHERE deleted_at IS NULL;

-- name: FindMovieByID :one
SELECT id, title, description, duration_minutes, rating FROM movies WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateMovieByID :execrows
UPDATE movies SET title = $1, description = $2, duration_minutes = $3, rating = $4 WHERE id = $5 AND deleted_at IS NULL;

-- name: DeleteMovieByID :execrows
UPDATE movies SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL;