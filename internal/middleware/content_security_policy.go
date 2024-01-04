package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

var (
	DefaultContentSecurityPolicy func(next http.Handler) http.Handler
)

func init() {
	DefaultContentSecurityPolicy = middleware.SetHeader("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; img-src 'self' data:; font-src 'self';")
}

func ContentSecurityPolicy(next http.Handler) http.Handler {
	return DefaultContentSecurityPolicy(next)
}
