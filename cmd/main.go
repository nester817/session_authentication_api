package main

import (
	"context"
	"fmt"
	"log"
	"os"

	storage "github.com/nester817/session_authentication_api.git/pkg/storage/postgres"
)

func main() {
	ctx := context.Background()

	dburl, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewDb(ctx, dburl, 5)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)
}

func GetConfig() (string, error) {
	url := os.Getenv("DB_URL")
	if url == "" {
		return "", fmt.Errorf("error getting config")
	}

	return url, nil
}
