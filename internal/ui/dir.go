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
	"github.com/seb-schulz/onegate/internal/ui/config"
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
	ctx, err := api.Context(config.DefaultBuildOptions(path.Join(path.Dir(fileName()), "_client"), path.Join(outdir(), "static"), false))
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

func Template(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFS(os.DirFS(path.Join(path.Dir(fileName()), "_templates")), "*.tmpl")
		if err != nil {
			log.Fatalln("Cannot read template files: ", err)
		}

		t.ExecuteTemplate(w, filename, fromContext(r.Context()))
	}
}
