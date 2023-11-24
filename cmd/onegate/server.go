package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/seb-schulz/onegate/graph"
	"github.com/seb-schulz/onegate/internal/ui"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PasskeyID string `gorm:"type:BLOB(16);default:RANDOM_BYTES(16);not null"`
}

type Passkeys struct {
	gorm.Model
	UserID        int
	User          User
	Username      string `gorm:"type:VARCHAR(255);not null"`
	PublicKeySpki []byte `gorm:"type:BLOB"`
	Backup        bool
}

var defaultConfig = []byte(`
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

	if err := db.AutoMigrate(User{}, Passkeys{}); err != nil {
		log.Fatalln("Migration failed: ", err)
	}

	// Manual migration was added because tags generated multiple indexes
	if !db.Migrator().HasIndex(&User{}, "idx_user_passkey_id_uniq") {
		db.Exec("CREATE UNIQUE INDEX idx_user_passkey_id_uniq ON users(passkey_id)")
	}

	http.Handle("/", ui.Template("index.html.tmpl"))
	http.Handle("/favicon.ico", ui.PublicFile())
	http.Handle("/robots.txt", ui.PublicFile())
	http.Handle("/hello", ui.Template("hello.html.tmpl"))
	http.Handle("/static/", ui.StaticFiles())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/query", srv)

	addGraphQLPlayground()

	port := viper.GetString("httpPort")
	fmt.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln(err)
	}
}
