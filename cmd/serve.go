package cmd

import (
	"fmt"
	"log"
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

func runServeCmd(cmd *cobra.Command, args []string) {
	db, err := gorm.Open(mysql.Open(config.Default.DB.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if config.Default.DB.Debug {
		db = db.Debug()
	}

	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: config.Default.RelyingParty.Name,
		RPID:          config.Default.RelyingParty.ID,
		RPOrigins:     config.Default.RelyingParty.Origins,
	})
	if err != nil {
		fmt.Println(err)
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

	port := config.Default.HttpPort
	fmt.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalln(err)
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	Run:   runServeCmd,
}
