package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/IMBoBx/url-shortener/internal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL env var is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8800"
		log.Println("defaulting to port 8800")
	}

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	createTable := `
		CREATE TABLE IF NOT EXISTS urls (
			original TEXT NOT NULL,
			short TEXT PRIMARY KEY
		);`

	_, err = dbpool.Exec(ctx, createTable)
	if err != nil {
		log.Fatal(err)
	}

	server := internal.NewServer(dbpool)

	log.Fatal(http.ListenAndServe(":"+port, server))
}
