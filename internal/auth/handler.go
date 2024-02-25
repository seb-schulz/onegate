package auth

import (
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

	// TODO: Replace will callback handler
	// route.With(usermgr.Middleware).Get("/callback", defaultAuthorizationResponseHandler.ServeHTTP)

	tokenHandler := &tokenHandler{
		clientByClientID:    clientByClientID,
		authorizationByCode: authorizationByCode,
	}
	route.Post("/token", tokenHandler.ServeHTTP)

	return route
}
