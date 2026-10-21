package movie

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateMovieRequest struct {
	Title           string          `json:"title" validate:"required"`
	Description     string          `json:"description" validate:"required"`
	DurationMinutes int             `json:"duration_mintes" validate:"required,min=1,max=500"`
	Rating          decimal.Decimal `json:"rating" validate:"required,rating"`
}

type UpdateMovieRequest struct {
	Title           string          `json:"title" validate:"required"`
	Description     string          `json:"description" validate:"required"`
	DurationMinutes int             `json:"duration_mintes" validate:"required,min=1,max=500"`
	Rating          decimal.Decimal `json:"rating" validate:"required,rating"`
}

type MovieResponseDTO struct {
	ID              uuid.UUID       `json:"id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	DurationMinutes int             `json:"duration_minutes"`
	Rating          decimal.Decimal `json:"rating"`
}
