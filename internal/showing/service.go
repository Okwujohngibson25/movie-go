package showing

import (
	"context"
	"time"

	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/google/uuid"
)

type ShowingRepository interface {
	GetAllShowings(ctx context.Context) ([]domains.Showing, error)
	CreateShowing(ctx context.Context, req domains.CreateShowing) error
	// HasConflictingShowing(ctx context.Context, screenID uuid.UUID, startTime time.Time, endTime time.Time) (bool, error)
}

type Movie interface {
	FindMovieByID(ctx context.Context, id uuid.UUID) (domains.Movie, error)
}

type ShowingService struct {
	repo  ShowingRepository
	movie Movie
}

func NewShowingService(repo ShowingRepository, movie Movie) *ShowingService {
	return &ShowingService{repo: repo, movie: movie}
}

func (s *ShowingService) GetAllShowings(ctx context.Context) ([]domains.Showing, error) {
	showings, err := s.repo.GetAllShowings(ctx)
	if err != nil {
		return nil, err
	}

	return showings, nil
}

func (s *ShowingService) CreateShowing(ctx context.Context, req CreateShowingInput) error {
	// find movie by ID
	movieData, err := s.movie.FindMovieByID(ctx, req.MovieID)
	if err != nil {
		return err
	}

	showingData := domains.CreateShowing{
		ID:        uuid.New(),
		MovieID:   req.MovieID,
		ScreenID:  req.ScreenID,
		StartTime: req.StartTime,
		EndTime:   req.StartTime.Add(time.Duration(movieData.DurationMinutes) * time.Minute),
	}

	if err := s.repo.CreateShowing(ctx, showingData); err != nil {
		return err
	}

	return nil
}
