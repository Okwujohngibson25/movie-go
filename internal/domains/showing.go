package domains

import (
	"time"

	"github.com/google/uuid"
)

type ShowingMovie struct {
	Title string `json:"title"`
}

type ShowingScreen struct {
	Name string `json:"name"`
}

type Showing struct {
	ID        uuid.UUID
	Movie     ShowingMovie
	Screen    ShowingScreen
	StartTime time.Time
	EndTime   time.Time
}

type CreateShowing struct {
	ID        uuid.UUID
	MovieID   uuid.UUID
	ScreenID  uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}
