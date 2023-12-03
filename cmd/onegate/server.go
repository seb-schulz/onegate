package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/seb-schulz/onegate/graph"
	"github.com/seb-schulz/onegate/internal/jwt"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/ui"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var defaultConfig = []byte(`
jwt:
  header: x-jwt-token
  secret: "NOT_CONFIGURED_YET"
  expires_in: 1h
  valid_methods: ["HS256", "HS384", "HS512"]
rp:
  name: "NOT_CONFIGURED_YET"
  id: "NOT_CONFIGURED_YET"
db:
  dsn: "NOT_CONFIGURED_YET"
httpPort: 9000
`)

func main() {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	db, err := gorm.Open(mysql.Open(viper.GetString("db.dsn")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(model.User{}, model.Passkeys{}, model.Session{}); err != nil {
		log.Fatalln("Migration failed: ", err)
	}

	// Manual migration was added because tags generated multiple indexes
	if !db.Migrator().HasIndex(&model.User{}, "idx_user_passkey_id_uniq") {
		db.Exec("CREATE UNIQUE INDEX idx_user_passkey_id_uniq ON users(passkey_id)")
	}

	http.Handle("/favicon.ico", ui.PublicFile())
	http.Handle("/robots.txt", ui.PublicFile())

	sessionMiddleware := middleware.SessionMiddleware(db)

	http.Handle("/", sessionMiddleware(ui.Template("index.html.tmpl", func() any {
		token, err := jwt.GenerateJwtToken(jwt.AnonymousUser)
		if err != nil {
			panic(err)
		}
		return map[string]any{"jwtInitToken": token, "jwtHeader": viper.GetString("jwt.header")}
	})))

	http.Handle("/static/", ui.StaticFiles())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db}}))
	http.Handle("/query", sessionMiddleware(srv))
	// http.Handle("/query", srv)

	addGraphQLPlayground()

	port := viper.GetString("httpPort")
	fmt.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln(err)
	}
}
