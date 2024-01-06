package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/fcgi"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/graph"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/ui"
	"github.com/seb-schulz/onegate/internal/utils"
)

type (
	ServeType int8

	RouterLimitConfig struct {
		RequestLimit int
		WindowLength time.Duration
	}

	RouterConfig struct {
		DbDebug  bool
		Webauthn webauthn.Config
		Limit    RouterLimitConfig
	}

	ServerConfig struct {
		Router    RouterConfig
		ServeType ServeType
		HttpPort  string
	}
)

const (
	ServeTypeHttp ServeType = iota
	ServeTypeFcgi
)

func newRouter(config *RouterConfig) (http.Handler, error) {
	db, err := utils.OpenDatabase(utils.WithDebugOption(config.DbDebug))
	if err != nil {
		return nil, err
	}

	webAuthn, err := webauthn.New(&config.Webauthn)
	if err != nil {
		return nil, fmt.Errorf("cannot configure WebAuth: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.ContentSecurityPolicy)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(httprate.LimitByRealIP(config.Limit.RequestLimit, config.Limit.WindowLength))
		r.Use(middleware.SessionMiddleware(db))

		r.Get("/login/{token}", middleware.LoginHandler(db))

		srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db, WebAuthn: webAuthn}}))

		r.Handle("/query", csrfMitigationMiddleware(srv))
		addGraphQLPlayground(r)
		r.Handle("/*", ui.Template("index.html.tmpl", func() any {
			return map[string]any{}
		}))
	})

	r.Group(func(r chi.Router) {
		r.Handle("/favicon.ico", ui.PublicFile())
		r.Handle("/robots.txt", ui.PublicFile())
		r.Mount("/static", ui.StaticFiles())
	})

	return r, nil
}

func Serve(config *ServerConfig) error {
	r, err := newRouter(&config.Router)
	if err != nil {
		return err
	}

	switch config.ServeType {
	case ServeTypeHttp:
		if config.HttpPort == "" {
			return fmt.Errorf("http port not defined")
		}

		port := config.HttpPort
		log.Println("Server listening on port ", port)
		if err := http.ListenAndServe(":"+port, r); err != nil {
			return fmt.Errorf("cannot run server: %v", err)
		}
	case ServeTypeFcgi:
		if err := fcgi.Serve(nil, r); err != nil {
			return fmt.Errorf("cannot run server: %v", err)
		}
	default:
		panic("cannot run any server type")
	}
	return nil
}
