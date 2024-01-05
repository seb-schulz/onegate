//go:build !production

package server

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

func addGraphQLPlayground(r chi.Router) {
	r.Get("/play", playground.HandlerWithHeaders("GraphQL playground", "/query", map[string]string{"X-Onegate-Csrf-Protection": "1"}, nil))
}
