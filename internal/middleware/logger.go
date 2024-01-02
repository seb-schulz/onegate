package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/seb-schulz/onegate/internal/config"
)

var (
	DefaultLogger func(next http.Handler) http.Handler
)

func init() {
	opt := httplog.Options{
		JSON:            true,
		LogLevel:        slog.LevelInfo,
		Concise:         true,
		RequestHeaders:  true,
		ResponseHeaders: true,
		TimeFieldFormat: time.RFC3339,
		Tags: map[string]string{
			"version": "latest",
		},
		QuietDownRoutes: []string{"/query"},
		QuietDownPeriod: time.Minute,
	}

	if config.Config.Logger.File != "" {
		f, err := os.OpenFile(config.Config.Logger.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0444)
		if err != nil {
			panic(err)
		}
		opt.Writer = f
	}
	DefaultLogger = httplog.RequestLogger(httplog.NewLogger("onegate", opt))
}

func Logger(next http.Handler) http.Handler {
	return DefaultLogger(next)
}
