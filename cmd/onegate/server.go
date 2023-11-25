package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/golang-jwt/jwt/v5"
	"github.com/seb-schulz/onegate/graph"
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
rp:
  name: "NOT_CONFIGURED_YET"
  id: "NOT_CONFIGURED_YET"
db:
  dsn: "NOT_CONFIGURED_YET"
httpPort: 9000
`)

func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(viper.GetString("jwt.header"))
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access forbidden"))
			return
		}
		log.Println("Request for graph", token)
		// TODO: Verify user

		next.ServeHTTP(w, r)
	})
}

func generateJwtToken(userID int) (string, error) {
	characters := "ABCDEFGHIJKLMOPQRSTUVWXYZabcdefghijklmopqrstuvwxyz0123456789"
	id_runes := make([]byte, 4)
	for i := range id_runes {
		id_runes[i] = characters[rand.Intn(len(characters))]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.expires_in"))),
		ID:        string(id_runes),
		Subject:   fmt.Sprintf("%x", userID),
	})

	secret, err := base64.StdEncoding.DecodeString(viper.GetString("jwt.secret"))
	if err != nil {
		return "", err
	}

	return token.SignedString(secret)
}

func main() {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	db, err := gorm.Open(mysql.Open(viper.GetString("db.dsn")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(model.User{}, model.Passkeys{}); err != nil {
		log.Fatalln("Migration failed: ", err)
	}

	// Manual migration was added because tags generated multiple indexes
	if !db.Migrator().HasIndex(&model.User{}, "idx_user_passkey_id_uniq") {
		db.Exec("CREATE UNIQUE INDEX idx_user_passkey_id_uniq ON users(passkey_id)")
	}

	http.Handle("/", ui.Template("index.html.tmpl", func() any {
		token, err := generateJwtToken(model.AnonymousUserID)
		if err != nil {
			panic(err)
		}
		return map[string]any{"jwtInitToken": token, "jwtHeader": viper.GetString("jwt.header")}
	}))

	http.Handle("/favicon.ico", ui.PublicFile())
	http.Handle("/robots.txt", ui.PublicFile())

	http.Handle("/static/", ui.StaticFiles())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/query", jwtAuthMiddleware(srv))

	addGraphQLPlayground()

	port := viper.GetString("httpPort")
	fmt.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln(err)
	}
}
