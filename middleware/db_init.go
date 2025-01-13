package gc_middleware

import (
	"context"

	"gc.yashk.dev/db"
	"gc.yashk.dev/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDb(context context.Context) (*db.Queries, *pgxpool.Pool, error) {
	dsn := env.DATABASE_URL
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, err
	}

	pool, err := pgxpool.NewWithConfig(context, config)
	if err != nil {
		return nil, nil, err
	}

	queries := db.New(pool)
	return queries, pool, nil
}
