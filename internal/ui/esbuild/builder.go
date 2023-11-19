package main

import (
	"log"
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

func main() {
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{"src/app.tsx"},
		Bundle:            true,
		AbsWorkingDir:     path.Join(path.Dir(fileName()), "../_client"),
		Outdir:            path.Join(path.Dir(fileName()), "../_build/static"),
		Write:             true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
	})
	if len(result.Errors) != 0 {
		for _, err := range result.Errors {
			log.Println(err)
		}
		os.Exit(1)
	}
}
