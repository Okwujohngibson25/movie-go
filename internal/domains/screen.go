package domains

import "github.com/google/uuid"

type Screen struct {
	ID       uuid.UUID
	Name     string
	Capacity int32
}
