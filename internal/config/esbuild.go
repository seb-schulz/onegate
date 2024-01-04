package config

import (
	"github.com/evanw/esbuild/pkg/api"
)

func DefaultBuildOptions(absWorkingDir, outdir string, prod bool) api.BuildOptions {
	opts := api.BuildOptions{
		EntryPoints:       []string{"src/app.tsx"},
		Bundle:            true,
		AbsWorkingDir:     absWorkingDir,
		Outdir:            outdir,
		Write:             true,
		MinifyWhitespace:  prod,
		MinifyIdentifiers: prod,
		MinifySyntax:      prod,
		Loader: map[string]api.Loader{
			".woff2": api.LoaderFile,
			".woff":  api.LoaderFile,
		},
	}

	if !prod {
		opts.Sourcemap = api.SourceMapInline
	}

	return opts
}
