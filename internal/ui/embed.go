//go:build embedded

//go:generate go run ./esbuild

package ui

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed _build _templates _public
var UI embed.FS

func StaticFiles() http.Handler {
	uiFS, err := fs.Sub(UI, "_build")
	if err != nil {
		log.Fatal("failed to get ui fs", err)
	}

	return http.FileServer(http.FS(uiFS))
}

func Template(filename string) http.Handler {
	tmplFs, err := fs.Sub(UI, "_templates")
	if err != nil {
		log.Fatal("failed to get template fs", err)
	}

	t, err := template.ParseFS(tmplFs, "*.tmpl")
	if err != nil {
		log.Fatalln("Cannot read template files: ", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t.ExecuteTemplate(w, filename, nil)
	})
}

func PublicFile() http.Handler {
	pFS, err := fs.Sub(UI, "_public")
	if err != nil {
		log.Fatal("failed to get ui fs", err)
	}

	return http.FileServer(http.FS(pFS))
}
