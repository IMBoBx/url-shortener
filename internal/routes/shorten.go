package routes

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type shortenReqBody struct {
	URL string `json:"url"`
}

type shortenResBody struct {
	ShortURL string `json:"shortUrl"`
}

func generateCode() (string, error) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 8)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		code[i] = alphabet[n.Int64()]
	}

	return string(code), nil
}

// POST /shorten
func ShortenHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		var bodyJson shortenReqBody
		err = json.Unmarshal(bodyBytes, &bodyJson)
		if err != nil {
			http.Error(w, "invalid JSON structure", http.StatusBadRequest)
			return
		}

		if bodyJson.URL == "" {
			http.Error(w, "url is required", http.StatusBadRequest)
			return
		}

		var code string
		var saved bool

		for range 50 {
			query := `INSERT INTO urls (original, short) VALUES($1, $2);`

			code, err = generateCode()
			if err != nil {
				http.Error(w, "failed to generate short code", http.StatusInternalServerError)
				return
			}

			_, err = dbpool.Exec(r.Context(), query, bodyJson.URL, code)

			if err == nil {
				saved = true
				break
			}
		}

		if !saved {
			http.Error(w, "failed to save short URL", http.StatusInternalServerError)
			return
		}

		scheme := "http"
		baseUrl := scheme + "://" + r.Host

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(shortenResBody{ShortURL: baseUrl + "/" + code})

		fmt.Printf("POST /shorten: Shortened %s to %s\n", bodyJson.URL, baseUrl+"/"+code)
	}
}
