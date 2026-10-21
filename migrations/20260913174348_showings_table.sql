-- +goose Up
-- +goose StatementBegin

--create showing table 
CREATE TABLE showings(
    id UUID PRIMARY KEY,
    screen_id UUID NOT NULL REFERENCES screens(id),
    movie_id UUID NOT NULL REFERENCES movies(id),
    start_time TIMESTAMPTZ NOT NULL,
    end_time   TIMESTAMPTZ NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS showings;
-- +goose StatementEnd
