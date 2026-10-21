package screen

import (
	"context"

	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/google/uuid"
)

type ScreenRepository interface {
	CreateScreenWithSeats(ctx context.Context, data domains.Screen) error
	GetAllScreens(ctx context.Context) ([]domains.Screen, error)
	GetScreensByID(ctx context.Context, Id uuid.UUID) (domains.Screen, error)
	UpdateScreensByID(ctx context.Context, data domains.Screen) error
	GetSeatsByScreenID(ctx context.Context, id uuid.UUID) ([]domains.Seats, error)
}

type ScreenService struct {
	repo ScreenRepository
}

func NewScreenService(repo ScreenRepository) *ScreenService {
	return &ScreenService{
		repo: repo,
	}
}

func (r *ScreenService) CreateScreen(ctx context.Context, req CreateScreenInput) error {
	screenData := domains.Screen{
		ID:       uuid.New(),
		Name:     req.Name,
		Capacity: int32(req.Capacity),
	}

	if err := r.repo.CreateScreenWithSeats(ctx, screenData); err != nil {
		return err
	}

	return nil
}

func (r *ScreenService) GetAllScreens(ctx context.Context) ([]domains.Screen, error) {
	screens, err := r.repo.GetAllScreens(ctx)
	if err != nil {
		return nil, err
	}

	return screens, nil
}

func (r *ScreenService) GetScreensByID(ctx context.Context, screenID uuid.UUID) (domains.Screen, error) {
	screen, err := r.repo.GetScreensByID(ctx, screenID)
	if err != nil {
		return domains.Screen{}, err
	}

	return screen, nil
}

func (r *ScreenService) UpdateScreensByID(ctx context.Context, screenID uuid.UUID, req UpdateScreenInput) error {
	ScreenData := domains.Screen{
		ID:   screenID,
		Name: req.Name,
	}

	if err := r.repo.UpdateScreensByID(ctx, ScreenData); err != nil {
		return err
	}

	return nil
}

func (r *ScreenService) GetSeatsByScreenID(ctx context.Context, screenID uuid.UUID) ([]domains.Seats, error) {
	seats, err := r.repo.GetSeatsByScreenID(ctx, screenID)
	if err != nil {
		return nil, err
	}

	return seats, nil
}
