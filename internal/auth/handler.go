package auth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func NewHandler() http.Handler {
	route := chi.NewRouter()

	authorizationRequestHandler := authorizationRequestHandler{
		clientByClientID:    clientByClientID,
		loginUrl:            url.URL{Path: "/login"},
		createAuthorization: createAuthorization,
	}
	route.Get("/auth", authorizationRequestHandler.ServeHTTP)

	callbackRedirectHandler := &callbackRedirectHandler{
		currentAuthorization: func(ctx context.Context) (authorization, error) {
			return firstAuthorization(ctx)
		},
	}
	route.Get("/callback", callbackRedirectHandler.ServeHTTP)

	tokenHandler := &tokenHandler{
		clientByClientID:    clientByClientID,
		authorizationByCode: authorizationByCode,
	}
	route.Post("/token", tokenHandler.ServeHTTP)

	return route
}
