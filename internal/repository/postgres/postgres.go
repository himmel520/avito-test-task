package postgres

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func New(cfg *config.PGConfig) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), cfg.Conn)
	if err != nil {
		return nil, err
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Postgres{DB: pool}, nil
}
