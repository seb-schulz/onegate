package server

import "net/http"

func csrfMitigation(next http.Handler) http.Handler {
	unauthorizedRequest := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized request"))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, ok := r.Header["X-Onegate-Csrf-Protection"]
		if !ok {
			unauthorizedRequest(w)
			return
		}

		if len(value) > 1 {
			unauthorizedRequest(w)
			return
		}

		if value[0] != "1" {
			unauthorizedRequest(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
