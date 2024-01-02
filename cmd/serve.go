package cmd

import (
	"fmt"
	"net/http"
	"net/http/fcgi"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/graph"
	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/ui"
	"github.com/seb-schulz/onegate/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

func runServeCmd(cmd *cobra.Command, args []string) error {
	db, err := utils.OpenDatabase(utils.WithDebugOption(config.Config.DB.Debug))
	if err != nil {
		return err
	}

	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: config.Config.RelyingParty.Name,
		RPID:          config.Config.RelyingParty.ID,
		RPOrigins:     config.Config.RelyingParty.Origins,
	})
	if err != nil {
		return fmt.Errorf("cannot configure WebAuth: %v", err)
	}

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(httprate.LimitByRealIP(config.Config.Server.Limit.RequestLimit, config.Config.Server.Limit.WindowLength))
		r.Use(middleware.SessionMiddleware(db))

		r.Get("/login/{token}", middleware.LoginHandler(db))

		srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db, WebAuthn: webAuthn}}))

		r.Handle("/query", srv)
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

	if config.Config.Server.Kind == config.ServerKindHttp {
		if config.Config.Server.HttpPort == "" {
			return fmt.Errorf("http port not defined")
		}

		port := config.Config.Server.HttpPort
		fmt.Println("Server listening on port ", port)
		if err := http.ListenAndServe(":"+port, r); err != nil {
			return fmt.Errorf("cannot run server: %v", err)
		}
	} else if config.Config.Server.Kind == config.ServerKindFcgi {
		if err := fcgi.Serve(nil, r); err != nil {
			return fmt.Errorf("cannot run server: %v", err)
		}
	}
	return fmt.Errorf("cannot run any server type")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	RunE:  runServeCmd,
}
