package showing

import (
	"time"

	"github.com/google/uuid"
)

type CreateShowingInput struct {
	MovieID   uuid.UUID
	ScreenID  uuid.UUID
	StartTime time.Time
}
