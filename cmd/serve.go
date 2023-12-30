package cmd

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/graph"
	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/ui"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

func runServeCmd(cmd *cobra.Command, args []string) error {
	db, err := gorm.Open(mysql.Open(config.Config.DB.Dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	if config.Config.DB.Debug {
		db = db.Debug()
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
		r.Use(chi_middleware.Logger)
		r.Use(chi_middleware.Recoverer)
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

	port := config.Config.HttpPort
	fmt.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		return fmt.Errorf("cannot run server: %v", err)
	}
	return nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	RunE:  runServeCmd,
}
