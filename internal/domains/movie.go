package domains

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Movie struct {
	ID              uuid.UUID
	Title           string
	Description     string
	DurationMinutes int
	Rating          decimal.Decimal
}
