-- +goose Up
-- +goose StatementBegin

--create screens table 

CREATE TABLE screens(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    capacity INTEGER NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS screens;
-- +goose StatementEnd
