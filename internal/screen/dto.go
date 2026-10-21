package screen

import "github.com/google/uuid"

type CreateScreenRequest struct {
	Name     string `json:"name" validate:"required"`
	Capacity int32  `json:"capacity" validate:"required,min=1,max=1000"`
}

type UpdateScreenRequest struct {
	Name string `json:"name" validate:"required"`
}

type ScreenResponseDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Capacity int32     `json:"capacity"`
}

type SeatsResponseDTO struct {
	ID         uuid.UUID `json:"id"`
	ScreenID   uuid.UUID `json:"screen_id"`
	SeatNumber int64     `json:"seat_number"`
}
