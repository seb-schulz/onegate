//go:build production

package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func addGraphQLPlayground(r chi.Router) {
	r.Get("/play", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 Not found", http.StatusNotFound)
	}))
}
