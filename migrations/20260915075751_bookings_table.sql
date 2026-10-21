-- +goose Up
-- +goose StatementBegin

--create bookings table
CREATE TABLE bookings(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    showing_id UUID NOT NULL REFERENCES showings(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
-- +goose StatementEnd
