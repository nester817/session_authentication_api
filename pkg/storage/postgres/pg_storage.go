package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	customer "github.com/nester817/session_authentication_api.git/pkg/storage"
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
			if _, err := client.Exec(
				ctx,
				`CREATE TABLE IF NOT EXISTS public.customer (
    				id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    				name TEXT NOT NULL,
    				email TEXT NOT NULL UNIQUE,
					password TEXT NOT NULL
				);`,
			); err != nil {
				return nil, err
			}

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

func (db *Db) GetUserByEmail(ctx context.Context, email, password string) (*customer.Customer, error) {
	var user customer.Customer

	if err := db.client.QueryRow(
		ctx,
		`SELECT id, name, email 
		FROM public.customer 
		WHERE email = $1 AND password = $2`,
		email, password,
	).Scan(&user.Id, &user.Name, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Db) InsertUser(ctx context.Context, user *customer.Customer) error {
	_, err := db.client.Exec(
		ctx,
		`INSERT INTO public.customer (name, email, password) 
		VALUES ($1, $2, $3)`,
		user.Name, user.Email, user.Password,
	)
	return err
}

func (db *Db) DeleteUserByEmail(ctx context.Context, email, password string) error {
	_, err := db.client.Exec(
		ctx,
		`DELETE FROM public.customer
		WHERE email = $1 AND password = $2`,
		email, password,
	)

	return err
}
