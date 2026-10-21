package movie

import "github.com/shopspring/decimal"

type MovieInput struct {
	Title           string
	Description     string
	DurationMinutes int
	Rating          decimal.Decimal
}
