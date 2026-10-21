package postgres

import (
	"context"
	"errors"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/database/db"
	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSqlcRepository struct {
	queries *db.Queries
}

func NewUserSqlcRepository(database *pgxpool.Pool) *UserSqlcRepository {
	return &UserSqlcRepository{
		queries: db.New(database),
	}
}

func (r *UserSqlcRepository) CreateUser(ctx context.Context, userData domains.User) error {
	CreateUserParams := db.CreateUserParams{
		ID: pgtype.UUID{
			Bytes: userData.ID,
			Valid: true,
		},
		Fullname:     userData.Fullname,
		Email:        userData.Email,
		Role:         db.AccountRole(userData.Role),
		PasswordHash: userData.PasswordHash,
	}
	err := r.queries.CreateUser(ctx, CreateUserParams)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return apperrors.ErrEmailAlreadyExist
			}
		}
		return err
	}
	return nil
}

func (r *UserSqlcRepository) FindByEmail(ctx context.Context, email string) (domains.User, error) {
	user, err := r.queries.FindByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return domains.User{}, apperrors.ErrEmailDoesNotExist
	}
	if err != nil {
		return domains.User{}, err
	}

	return domains.User{
		ID:           user.ID.Bytes,
		Fullname:     user.Fullname,
		Email:        user.Email,
		Role:         string(user.Role),
		PasswordHash: user.PasswordHash,
	}, nil
}
