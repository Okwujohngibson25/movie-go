-- name: CreateShowing :exec
INSERT INTO showings (id, screen_id, movie_id, start_time, end_time) VALUES($1, $2, $3, $4, $5);

-- name: GetAllShowings :many 
SELECT
    s.id,
    s.screen_id,
    s.start_time,
    s.end_time,
    m.title AS movie_title,
    sc.title AS screen_title
FROM showings AS s
JOIN movies AS m ON m.id = s.movie_id
JOIN screens AS sc ON sc.id = s.screen_id;


