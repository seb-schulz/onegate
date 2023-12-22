//go:build !production

package cmd

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
)

func addGraphQLPlayground() {
	http.Handle("/play", playground.Handler("GraphQL playground", "/query"))

}
