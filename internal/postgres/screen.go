package postgres

import (
	"context"
	"errors"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/database/db"
	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScreenSqlcRepository struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewScreenSqlcRepository(con *pgxpool.Pool) *ScreenSqlcRepository {
	return &ScreenSqlcRepository{
		queries: db.New(con),
		db:      con,
	}
}

func (r *ScreenSqlcRepository) CreateScreenWithSeats(ctx context.Context, data domains.Screen) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	txQueries := r.queries.WithTx(tx)

	arg := db.CreateScreenParams{
		ID: pgtype.UUID{
			Bytes: data.ID,
			Valid: true,
		},
		Title:    data.Name,
		Capacity: int32(data.Capacity),
	}

	if err := txQueries.CreateScreen(ctx, arg); err != nil {
		return err
	}

	for seatNumber := 1; seatNumber <= int(data.Capacity); seatNumber++ {
		seatID := uuid.New()
		seatArg := db.CreateSeatParams{
			ID: pgtype.UUID{
				Bytes: seatID,
				Valid: true,
			},
			SeatNumber: int32(seatNumber),
			ScreenID: pgtype.UUID{
				Bytes: data.ID,
				Valid: true,
			},
		}

		if err := txQueries.CreateSeat(ctx, seatArg); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ScreenSqlcRepository) GetAllScreens(ctx context.Context) ([]domains.Screen, error) {
	screen, err := r.queries.GetAllScreens(ctx)
	if err != nil {
		return nil, err
	}

	screenData := make([]domains.Screen, 0, len(screen))

	for _, value := range screen {
		screenData = append(screenData, domains.Screen{
			ID:       value.ID.Bytes,
			Name:     value.Title,
			Capacity: value.Capacity,
		})
	}

	return screenData, nil
}

func (r *ScreenSqlcRepository) GetScreensByID(ctx context.Context, id uuid.UUID) (domains.Screen, error) {
	pgUUID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
	screenData, err := r.queries.GetScreensByID(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Screen{}, apperrors.ErrScreenNotFound
		}
		return domains.Screen{}, err
	}

	return domains.Screen{
		ID:       screenData.ID.Bytes,
		Name:     screenData.Title,
		Capacity: screenData.Capacity,
	}, nil
}

func (r *ScreenSqlcRepository) UpdateScreensByID(ctx context.Context, data domains.Screen) error {
	pgUUID := pgtype.UUID{
		Bytes: data.ID,
		Valid: true,
	}

	param := db.UpdateScreensByIDParams{
		ID:    pgUUID,
		Title: data.Name,
	}

	rows, err := r.queries.UpdateScreensByID(ctx, param)
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.ErrScreenNotFound
	}

	return nil
}

func (r *ScreenSqlcRepository) GetSeatsByScreenID(ctx context.Context, id uuid.UUID) ([]domains.Seats, error) {
	pgUUID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
	seats, err := r.queries.GetSeatsByScreenID(ctx, pgUUID)
	if err != nil {
		return nil, err
	}

	seatData := make([]domains.Seats, 0, len(seats))
	for _, value := range seats {
		seatData = append(seatData, domains.Seats{
			ID:         value.ID.Bytes,
			ScreenID:   value.ScreenID.Bytes,
			SeatNumber: int64(value.SeatNumber),
		})
	}

	return seatData, nil
}
