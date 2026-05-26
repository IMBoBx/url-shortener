package main

import (
	"context"
	"log"
	"net/http"

	"github.com/IMBoBx/url-shortener/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	server := internal.NewServer(dbpool)

	log.Fatal(http.ListenAndServe(":8800", server))
}
