package postgres

import (
	"context"

	"github.com/Okwujohngibson25/movie-go/internal/database/db"
	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShowingSqlcRepository struct {
	queries *db.Queries
}

func NewShowingSqlcRepository(pool *pgxpool.Pool) *ShowingSqlcRepository {
	return &ShowingSqlcRepository{
		queries: db.New(pool),
	}
}

func (r *ShowingSqlcRepository) CreateShowing(ctx context.Context, req domains.CreateShowing) error {
	id := pgtype.UUID{
		Bytes: req.ID,
		Valid: true,
	}

	screenID := pgtype.UUID{
		Bytes: req.ScreenID,
		Valid: true,
	}

	movieID := pgtype.UUID{
		Bytes: req.MovieID,
		Valid: true,
	}

	arg := db.CreateShowingParams{
		ID:        id,
		ScreenID:  screenID,
		MovieID:   movieID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := r.queries.CreateShowing(ctx, arg); err != nil {
		return err
	}

	return nil
}

func (r *ShowingSqlcRepository) GetAllShowings(ctx context.Context) ([]domains.Showing, error) {
	showings, err := r.queries.GetAllShowings(ctx)
	if err != nil {
		return nil, err
	}

	showingsData := make([]domains.Showing, 0, len(showings))

	for _, value := range showings {
		showingsData = append(showingsData, domains.Showing{
			ID: value.ID.Bytes,
			Movie: domains.ShowingMovie{
				Title: value.MovieTitle,
			},
			Screen: domains.ShowingScreen{
				Name: value.ScreenTitle,
			},
			StartTime: value.StartTime,
			EndTime:   value.EndTime,
		})
	}

	return showingsData, nil
}

// func (r *SqlcRepository) HasConflictingShowing(ctx context.Context, screenID uuid.UUID, startTime time.Time, endTime time.Time) (bool, error) {
// }
