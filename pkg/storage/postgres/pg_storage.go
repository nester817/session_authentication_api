package storagePostgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	customer "github.com/nester817/session_authentication_api.git/pkg/storage"
)

type DbPostgres struct {
	client *pgxpool.Pool
}

func NewDbPostgres(ctx context.Context, addr string, attempt int) (*DbPostgres, error) {
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
    				email TEXT NOT NULL UNIQUE,
					password TEXT NOT NULL
				);`,
			); err != nil {
				return nil, err
			}

			return &DbPostgres{
				client: client,
			}, nil
		} else {
			client.Close()
		}

		time.Sleep(1 * time.Second)
	}

	return nil, err
}

func (db *DbPostgres) GetCustomer(ctx context.Context, email string) (*customer.Customer, error) {
	var user customer.Customer

	if err := db.client.QueryRow(
		ctx,
		`SELECT id, email, password
		FROM public.customer 
		WHERE email = $1`,
		email,
	).Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DbPostgres) GetCustomerById(ctx context.Context, id string) (*customer.Customer, error) {
	var customer customer.Customer

	if err := db.client.QueryRow(ctx, `
		SELECT email
		FROM public.customer 
		WHERE id = $1
	`, id).Scan(&customer.Email); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (db *DbPostgres) InsertCustomer(ctx context.Context, user *customer.Customer) error {
	_, err := db.client.Exec(
		ctx,
		`INSERT INTO public.customer (email, password)
		VALUES ($1, $2)`,
		user.Email, &user.Password,
	)
	return err
}
