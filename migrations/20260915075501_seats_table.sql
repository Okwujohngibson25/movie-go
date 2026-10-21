-- +goose Up
-- +goose StatementBegin

--create seats table 
CREATE TABLE seats(
    id UUID PRIMARY KEY,
    seat_number INTEGER NOT NULL,
    screen_id UUID NOT NULL REFERENCES screens(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS seats;
-- +goose StatementEnd
