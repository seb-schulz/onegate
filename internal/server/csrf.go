package server

import (
	"net/http"
)

func csrfMitigationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, ok := r.Header["X-Onegate-Csrf-Protection"]
		if !ok {
			http.Error(w, "unauthorized request", http.StatusUnauthorized)
			return
		}

		if len(value) > 1 {
			http.Error(w, "unauthorized request", http.StatusUnauthorized)
			return
		}

		if value[0] != "1" {
			http.Error(w, "unauthorized request", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
