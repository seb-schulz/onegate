//go:build production

package main

import (
	"net/http"
)

func addGraphQLPlayground() {
	http.Handle("/play", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 Not found", http.StatusNotFound)
	}))
}
