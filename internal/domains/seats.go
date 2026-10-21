package domains

import "github.com/google/uuid"

type Seats struct {
	ID         uuid.UUID
	ScreenID   uuid.UUID
	SeatNumber int64
}
