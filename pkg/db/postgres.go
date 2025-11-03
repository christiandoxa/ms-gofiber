package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	apmpgx "go.elastic.co/apm/module/apmpgxv5/v2"

	"ms-gofiber/internal/config"
)

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(cfg.PGUrl)
	if err != nil {
		return nil, err
	}
	// APM instrumentation
	apmpgx.Instrument(conf.ConnConfig)

	conf.MaxConns = cfg.PGMaxConn
	conf.MinConns = cfg.PGMinConn
	conf.MaxConnLifetime = 60 * time.Minute
	conf.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}
	// ping
	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
