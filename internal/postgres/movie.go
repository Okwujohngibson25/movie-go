package postgres

import (
	"context"
	"errors"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/database/db"
	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/Okwujohngibson25/movie-go/internal/helper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieSqlcRepository struct {
	queries *db.Queries
}

func NewMovieSqlcRepository(conn *pgxpool.Pool) *MovieSqlcRepository {
	return &MovieSqlcRepository{
		queries: db.New(conn),
	}
}

func (r *MovieSqlcRepository) FindMovieByID(ctx context.Context, id uuid.UUID) (domains.Movie, error) {
	// convert uuid to pg.uuid
	pgID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
	movieData, err := r.queries.FindMovieByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Movie{}, apperrors.ErrMovieNotFound
		}
		return domains.Movie{}, err
	}

	rating, err := helper.NumericToDecimal(movieData.Rating)

	return domains.Movie{
		ID:              movieData.ID.Bytes,
		Title:           movieData.Title,
		Description:     movieData.Description,
		DurationMinutes: int(movieData.DurationMinutes),
		Rating:          rating,
	}, nil
}

func (r *MovieSqlcRepository) GetAllMovies(ctx context.Context) ([]domains.Movie, error) {
	movies, err := r.queries.GetAllMovies(ctx)
	if err != nil {
		return nil, err
	}

	movieData := make([]domains.Movie, 0, len(movies))

	for _, value := range movies {
		decRating, err := helper.NumericToDecimal(value.Rating)
		if err != nil {
			return nil, err
		}
		movieData = append(movieData, domains.Movie{
			ID:              value.ID.Bytes,
			Title:           value.Title,
			Description:     value.Description,
			DurationMinutes: int(value.DurationMinutes),
			Rating:          decRating,
		})
	}

	return movieData, nil
}

func (r *MovieSqlcRepository) CreateMovie(ctx context.Context, req domains.Movie) error {
	var rating pgtype.Numeric

	if err := rating.Scan(req.Rating.String()); err != nil {
		return err
	}

	arg := db.CreateMovieParams{
		ID: pgtype.UUID{
			Bytes: req.ID,
			Valid: true,
		},
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: int32(req.DurationMinutes),
		Rating:          rating,
	}

	_, err := r.queries.CreateMovie(ctx, arg)
	if err != nil {
		// ON CONFLICT DO NOTHING returns no rows when the movie already exists
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrMovieAlreadyExists
		}
		return err
	}

	return nil
}

func (r *MovieSqlcRepository) UpdateMovieByID(ctx context.Context, movieID uuid.UUID, req domains.Movie) error {
	var rating pgtype.Numeric

	if err := rating.Scan(req.Rating.String()); err != nil {
		return err
	}

	arg := db.UpdateMovieByIDParams{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: int32(req.DurationMinutes),
		Rating:          rating,
		ID: pgtype.UUID{
			Bytes: movieID,
			Valid: true,
		},
	}
	rows, err := r.queries.UpdateMovieByID(ctx, arg)
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.ErrMovieNotFound
	}

	return nil
}

func (r *MovieSqlcRepository) DeleteMovieByID(ctx context.Context, movieID uuid.UUID) error {
	rows, err := r.queries.DeleteMovieByID(ctx, pgtype.UUID{
		Bytes: movieID,
		Valid: true,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.ErrMovieNotFound
	}

	return nil
}
