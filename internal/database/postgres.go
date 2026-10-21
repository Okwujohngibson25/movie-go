package database

import (
	"context"
	"fmt"

	"github.com/Okwujohngibson25/movie-go/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	fmt.Println("successfully connected to postgres")

	return pool, nil
}
