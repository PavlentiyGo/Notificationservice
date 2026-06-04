package repository

import "github.com/jackc/pgx/v5/pgxpool"

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(
	pool *pgxpool.Pool,
) *ProfileRepository {
	return &ProfileRepository{
		db: pool,
	}
}
