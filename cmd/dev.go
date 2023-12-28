//go:build !production

package cmd

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

func addGraphQLPlayground(r chi.Router) {
	r.Get("/play", playground.Handler("GraphQL playground", "/query"))
}
