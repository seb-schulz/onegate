//go:build !embedded

package ui

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/evanw/esbuild/pkg/api"
)

func fileName() string {
	_, fileName, _, ok := runtime.Caller(1)
	if ok {
		return fileName
	}
	return ""
}

func outdir() string {
	return path.Join(path.Dir(fileName()), "_build")
}

func init() {
	ctx, err := api.Context(api.BuildOptions{
		EntryPoints:       []string{"src/app.tsx"},
		Bundle:            true,
		AbsWorkingDir:     path.Join(path.Dir(fileName()), "_client"),
		Outdir:            path.Join(outdir(), "static"),
		Write:             true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
	})
	if err != nil {
		log.Fatalln("Cannot init esbuild watch: ", err)
	}

	if err := ctx.Watch(api.WatchOptions{}); err != nil {
		log.Fatalln("Cannot watch ui files: ", err)
	}
}

func StaticFiles() http.Handler {
	return http.FileServer(http.Dir(outdir()))
}

func PublicFile() http.Handler {
	return http.FileServer(http.Dir(path.Join(path.Dir(fileName()), "_public")))
}

func Template(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFS(os.DirFS(path.Join(path.Dir(fileName()), "_templates")), "*.tmpl")
		if err != nil {
			log.Fatalln("Cannot read template files: ", err)
		}

		t.ExecuteTemplate(w, filename, nil)
	})
}
