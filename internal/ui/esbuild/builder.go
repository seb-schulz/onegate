package main

import (
	"log"
	"os"
	"path"
	"runtime"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/seb-schulz/onegate/internal/config"
)

func fileName() string {
	_, fileName, _, ok := runtime.Caller(1)
	if ok {
		return fileName
	}
	return ""
}

func main() {
	result := api.Build(config.DefaultBuildOptions(path.Join(path.Dir(fileName()), "../_client"), path.Join(path.Dir(fileName()), "../_build/static"), true))

	if len(result.Errors) != 0 {
		for _, err := range result.Errors {
			log.Println(err)
		}
		os.Exit(1)
	}
}
