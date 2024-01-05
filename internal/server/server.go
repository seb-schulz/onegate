package server

import (
	"fmt"
	"net/http"
	"net/http/fcgi"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/graph"
	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/ui"
	"github.com/seb-schulz/onegate/internal/utils"
)

type (
	routerLimitConfig struct {
		requestLimit int
		windowLength time.Duration
	}

	routerConfig struct {
		dbDebug  bool
		webauthn webauthn.Config
		limit    routerLimitConfig
	}
)

func newRouter(config *routerConfig) (http.Handler, error) {
	db, err := utils.OpenDatabase(utils.WithDebugOption(config.dbDebug))
	if err != nil {
		return nil, err
	}

	webAuthn, err := webauthn.New(&config.webauthn)
	if err != nil {
		return nil, fmt.Errorf("cannot configure WebAuth: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.ContentSecurityPolicy)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(httprate.LimitByRealIP(config.limit.requestLimit, config.limit.windowLength))
		r.Use(middleware.SessionMiddleware(db))

		r.Get("/login/{token}", middleware.LoginHandler(db))

		srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db, WebAuthn: webAuthn}}))

		r.Handle("/query", csrfMitigation(srv))
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

func Serve() error {
	r, err := newRouter(&routerConfig{
		dbDebug: config.Config.DB.Debug,
		webauthn: webauthn.Config{
			RPDisplayName: config.Config.RelyingParty.Name,
			RPID:          config.Config.RelyingParty.ID,
			RPOrigins:     config.Config.RelyingParty.Origins,
		},
		limit: routerLimitConfig{
			config.Config.Server.Limit.RequestLimit, config.Config.Server.Limit.WindowLength,
		},
	})
	if err != nil {
		return err
	}

	switch config.Config.Server.Kind {
	case config.ServerKindHttp:
		if config.Config.Server.HttpPort == "" {
			return fmt.Errorf("http port not defined")
		}

		port := config.Config.Server.HttpPort
		fmt.Println("Server listening on port ", port)
		if err := http.ListenAndServe(":"+port, r); err != nil {
			return fmt.Errorf("cannot run server: %v", err)
		}
	case config.ServerKindFcgi:
		if err := fcgi.Serve(nil, r); err != nil {
			return fmt.Errorf("cannot run server: %v", err)
		}
	default:
		panic("cannot run any server type")
	}
	return nil
}
