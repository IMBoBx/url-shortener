package routes

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RedirectHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		query := `
			SELECT original
			FROM urls
			WHERE short=$1;
		`

		var original string
		err := dbpool.QueryRow(r.Context(), query, code).Scan(&original)
		if err != nil {
			http.Error(w, "short URL not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, original, http.StatusFound)
		fmt.Printf("GET /%s: Redirected to %s\n", code, original)
	}
}
