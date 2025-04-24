package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	client *pgxpool.Pool
}

func NewDb(ctx context.Context, addr string, attempt int) (*Db, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 15
	config.MinConns = 5
	config.MaxConnLifetime = 120 * time.Minute
	config.MaxConnIdleTime = 120 * time.Minute

	var client *pgxpool.Pool

	for i := 0; i < attempt; i++ {
		client, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := client.Ping(ctx); err == nil {
			return &Db{
				client: client,
			}, nil
		} else {
			client.Close()
		}

		time.Sleep(1 * time.Second)
	}

	return nil, err
}
