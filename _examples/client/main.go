package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const tpl = `<!doctype html>
<html><head>
	<meta charset="utf-8">
	<title>oAuth2 Example</title>
  </head>
  <body>
	<p><a href="{{.Url}}">start oAuth</a></p>
	<p><code>{{.Url}}</code></p>
  </body>
</html>`

var verifier string

func index(conf *oauth2.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifier = oauth2.GenerateVerifier()

		url := conf.AuthCodeURL("state", oauth2.S256ChallengeOption(verifier))

		t, err := template.New("webpage").Parse(tpl)
		if err != nil {
			log.Fatalf("cannot parse template: %w", err)
		}
		t.Execute(w, map[string]any{
			"Url": url,
		})
		fmt.Fprint(w)
	})
}

func callback(conf *oauth2.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")

		tok, err := conf.Exchange(r.Context(), code, oauth2.VerifierOption(verifier))
		if err != nil {
			log.Fatal(err)
		}

		client := conf.Client(r.Context(), tok)

		state := r.FormValue("state")
		log.Println(state, client)
		fmt.Fprintf(w, "Hello World with token: %s", tok)
	})
}

func main() {
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./_examples/client/")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("EXAMPLE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	conf := &oauth2.Config{
		ClientID:     viper.GetString("Client.ID"),
		ClientSecret: viper.GetString("Client.secret"),
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  viper.GetString("Endpoint.AuthURL"),
			TokenURL: viper.GetString("Endpoint.TokenURL"),
		},
	}

	route := chi.NewRouter()
	route.Get("/cb", callback(conf))
	route.Get("/", index(conf))

	port := "9010"

	log.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, route); err != nil {
		log.Fatalf("cannot run server: %v", err)
	}
}
