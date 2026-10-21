package postgres

import (
	"testing"

	"github.com/Okwujohngibson25/movie-go/internal/database/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func TestMovieFromRow_InvalidRatingReturnsError(t *testing.T) {
	row := db.FindMovieByIDRow{
		ID:              pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Title:           "Inception",
		Description:     "dreams within dreams",
		DurationMinutes: 148,
		Rating:          pgtype.Numeric{Valid: false},
	}

	if _, err := movieFromRow(row); err == nil {
		t.Fatal("want error for invalid rating, got nil")
	}
}

func TestMovieFromRow_MapsAllFields(t *testing.T) {
	id := uuid.New()
	var rating pgtype.Numeric
	if err := rating.Scan("4.50"); err != nil {
		t.Fatalf("scan rating: %v", err)
	}
	row := db.FindMovieByIDRow{
		ID:              pgtype.UUID{Bytes: id, Valid: true},
		Title:           "Inception",
		Description:     "dreams within dreams",
		DurationMinutes: 148,
		Rating:          rating,
	}

	movie, err := movieFromRow(row)
	if err != nil {
		t.Fatalf("movieFromRow: %v", err)
	}

	if movie.ID != id || movie.Title != "Inception" || movie.Description != "dreams within dreams" || movie.DurationMinutes != 148 {
		t.Fatalf("movie = %+v, fields not mapped from row", movie)
	}
	if want := decimal.RequireFromString("4.50"); !movie.Rating.Equal(want) {
		t.Fatalf("rating = %s, want %s", movie.Rating, want)
	}
}
