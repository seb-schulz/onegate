package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
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

	http.Handle("/favicon.ico", ui.PublicFile())
	http.Handle("/robots.txt", ui.PublicFile())

	sessionMiddleware := middleware.SessionMiddleware(db)

	http.Handle("/", sessionMiddleware(ui.Template("index.html.tmpl", func() any {
		return map[string]any{}
	})))

	http.Handle("/static/", ui.StaticFiles())

	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: config.Default.RelyingParty.Name,
		RPID:          config.Default.RelyingParty.ID,
		RPOrigins:     config.Default.RelyingParty.Origins,
	})
	if err != nil {
		fmt.Println(err)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db, WebAuthn: webAuthn}}))
	http.Handle("/query", sessionMiddleware(srv))
	// http.Handle("/query", srv)

	addGraphQLPlayground()

	port := config.Default.HttpPort
	fmt.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln(err)
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	Run:   runServeCmd,
}
