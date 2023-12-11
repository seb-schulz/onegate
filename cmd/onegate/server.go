package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/graph"
	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/ui"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open(config.Default.DB.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if config.Default.DB.Debug {
		db = db.Debug()
	}

	if err := db.AutoMigrate(model.User{}, model.Credential{}, model.Session{}, model.AuthSession{}); err != nil {
		log.Fatalln("Migration failed: ", err)
	}

	// Manual migration was added because tags generated multiple indexes
	if !db.Migrator().HasIndex(&model.User{}, "idx_user_authn_id_uniq") {
		db.Exec("CREATE UNIQUE INDEX idx_user_authn_id_uniq ON users(authn_id)")
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
