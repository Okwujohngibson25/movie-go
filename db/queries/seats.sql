-- name: CreateSeat :exec
INSERT INTO seats (id, seat_number, screen_id) VALUES($1, $2, $3);

-- name: GetSeatsByScreenID :many
SELECT id, seat_number, screen_id FROM seats WHERE screen_id = $1;