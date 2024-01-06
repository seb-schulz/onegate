//go:build gqlgen

package graph

//go:generate go run github.com/99designs/gqlgen generate
//go:generate /bin/bash -c "(cd $(pwd)/../internal/ui/_client/ && npm run compile)"
