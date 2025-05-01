package main

import (
	"context"
	"fmt"
	"log"
	"os"

	router "github.com/nester817/session_authentication_api.git/pkg/api/gin_router"
	"github.com/nester817/session_authentication_api.git/pkg/api/server"
	storagePostgres "github.com/nester817/session_authentication_api.git/pkg/storage/postgres"
	storageRedis "github.com/nester817/session_authentication_api.git/pkg/storage/redis"
)

func main() {
	ctx := context.Background()

	pgurl, err := GetPostgresConfig()
	if err != nil {
		log.Fatal(err)
	}
	pgdb, err := storagePostgres.NewDbPostgres(ctx, pgurl, 5)
	if err != nil {
		log.Fatal(err)
	}

	redisurl, err := GetRedisConfig()
	if err != nil {
		log.Fatal(err)
	}
	redis := storageRedis.NewDbRedis(redisurl)

	app := router.NewHandler(pgdb, redis)
	httpServer := server.NewServer(":8080", app.InitRouter())

	if err := httpServer.Run(); err != nil {
		log.Fatal(err)
	}
}

func GetPostgresConfig() (string, error) {
	url := os.Getenv("DB_URL")
	if url == "" {
		return "", fmt.Errorf("error getting config")
	}

	return url, nil
}

func GetRedisConfig() (string, error) {
	url := os.Getenv("CACHE_URL")
	if url == "" {
		return "", fmt.Errorf("error getting config")
	}

	return url, nil
}
