package server

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func contentSecurityPolicyMiddleware(next http.Handler) http.Handler {
	middleware := middleware.SetHeader("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; img-src 'self' data:; font-src 'self';")
	return middleware(next)
}
