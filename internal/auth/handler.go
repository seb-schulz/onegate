package auth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/seb-schulz/onegate/internal/usermgr"
)

type Config struct {
	IssuerUrl  string
	PrivateKey *ecdsa.PrivateKey
}

func NewHandler(c *Config) http.Handler {
	route := chi.NewRouter()

	authorizationRequestHandler := authorizationRequestHandler{
		clientByClientID:    clientByClientID,
		loginUrl:            url.URL{Path: "/login"},
		createAuthorization: createAuthorization,
	}
	route.Get("/auth", authorizationRequestHandler.ServeHTTP)

	callbackRedirectHandler := &callbackRedirectHandler{
		currentAuthorization: func(ctx context.Context) (authorization, error) {
			return FirstAuthorization(ctx)
		},
		currentUser: usermgr.FromContext,
	}
	route.With(usermgr.Middleware).Get("/callback", callbackRedirectHandler.ServeHTTP)

	tokenHandler := &tokenHandler{
		issuerUrl:           c.IssuerUrl,
		privateKey:          c.PrivateKey,
		clientByClientID:    clientByClientID,
		authorizationByCode: authorizationByCode,
		deleteAuthorization: func(ctx context.Context, a authorization) error {
			return a.Delete(ctx)
		},
	}
	route.Post("/token", tokenHandler.ServeHTTP)

	return route
}

func RedirectWhenLoggedInAndAssigned(callbackURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authReq, err := FirstAuthorization(r.Context())
			if authReq == nil {
				next.ServeHTTP(w, r)
				return
			}

			if err != nil {
				// TODO: Evaluate proper error page
				slog.Warn(fmt.Sprintf("oAuth2 callback failed with: %v", err))
				http.NotFound(w, r)
				return
			}

			if authReq.InternalUserID != nil {
				http.Redirect(w, r, callbackURL, http.StatusSeeOther)
				return
			}

			user := usermgr.FromContext(r.Context())
			if user != nil {
				http.Redirect(w, r, callbackURL, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
