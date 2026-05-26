package routes

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RedirectHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		fmt.Print(code)
	}
}
