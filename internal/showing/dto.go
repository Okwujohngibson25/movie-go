package showing

import (
	"time"

	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/google/uuid"
)

type CreateShowingDTO struct {
	MovieID   string    `json:"movie_id"`
	ScreenID  string    `json:"screen_id"`
	StartTime time.Time `json:"start_time"`
}

type ShowingResponseDTO struct {
	ID        uuid.UUID             `json:"id"`
	Movie     domains.ShowingMovie  `json:"movie"`
	Screen    domains.ShowingScreen `json:"screen"`
	StartTime time.Time             `json:"start_time"`
	EndTime   time.Time             `json:"end_time"`
}
