package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/seb-schulz/onegate/graph"
	"github.com/seb-schulz/onegate/internal/ui"
)

const defaultPort = "9000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", ui.Template("index.html.tmpl"))
	http.Handle("/favicon.ico", ui.PublicFile())
	http.Handle("/robots.txt", ui.PublicFile())
	http.Handle("/hello", ui.Template("hello.html.tmpl"))
	http.Handle("/static/", ui.StaticFiles())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/query", srv)

	addGraphQLPlayground()

	fmt.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln(err)
	}
}
