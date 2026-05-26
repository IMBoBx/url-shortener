package internal

import (
	"net/http"

	"github.com/IMBoBx/url-shortener/internal/routes"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type Handler struct {
// 	path string;
// 	handlerFunc func(http.ResponseWriter, *http.Request);
// }

func NewServer(dbpool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{code}", routes.RedirectHandler(dbpool))
	mux.HandleFunc("POST /shorten", routes.ShortenHandler(dbpool))

	return mux
}
