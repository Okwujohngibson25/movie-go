package movie

import (
	"context"

	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/google/uuid"
)

type MovieRepository interface {
	FindMovieByID(ctx context.Context, id uuid.UUID) (domains.Movie, error)
	GetAllMovies(ctx context.Context) ([]domains.Movie, error)
	CreateMovie(ctx context.Context, req domains.Movie) error
	UpdateMovieByID(ctx context.Context, movieID uuid.UUID, req domains.Movie) error
	DeleteMovieByID(ctx context.Context, movieID uuid.UUID) error
}

type MovieService struct {
	repo MovieRepository
}

func NewMovieService(repo MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (r *MovieService) GetAllMovies(ctx context.Context) ([]domains.Movie, error) {
	movies, err := r.repo.GetAllMovies(ctx)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieService) GetMovieByID(ctx context.Context, movieID uuid.UUID) (movieData domains.Movie, err error) {
	// find and return movie + error
	movie, err := r.repo.FindMovieByID(ctx, movieID)
	if err != nil {
		return domains.Movie{}, err
	}

	return movie, nil
}

func (r *MovieService) CreateMovie(ctx context.Context, req MovieInput) error {
	movieData := domains.Movie{
		ID:              uuid.New(),
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Rating:          req.Rating,
	}

	if err := r.repo.CreateMovie(ctx, movieData); err != nil {
		return err
	}
	return nil
}

func (r *MovieService) UpdateMovieByID(ctx context.Context, movieID uuid.UUID, req MovieInput) error {
	movieData := domains.Movie{
		ID:              movieID,
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Rating:          req.Rating,
	}
	// update record
	if err := r.repo.UpdateMovieByID(ctx, movieID, movieData); err != nil {
		return err
	}

	return nil
}

func (r *MovieService) DeleteMovieByID(ctx context.Context, movieID uuid.UUID) error {
	if err := r.repo.DeleteMovieByID(ctx, movieID); err != nil {
		return err
	}

	return nil
}
