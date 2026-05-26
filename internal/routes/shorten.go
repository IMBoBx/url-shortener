package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ShortenHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
