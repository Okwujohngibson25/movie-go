-- +goose Up
-- +goose StatementBegin

--create booking seats
CREATE TABLE booking_seats(
    id UUID PRIMARY KEY,
    booking_id UUID NOT NULL REFERENCES bookings(id),
    seat_id UUID NOT NULL REFERENCES seats(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS booking_seats;
-- +goose StatementEnd
