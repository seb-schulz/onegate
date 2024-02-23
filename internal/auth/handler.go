package auth

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
	"github.com/seb-schulz/onegate/internal/usermgr"
)

var (
	defaultAuthorizationMgr             authorizationMgr
	defaultAuthorizationResponseHandler authorizationResponseHandler
)

func init() {
	defaultAuthorizationMgr = authorizationMgr{
		StorageManager: sessionmgr.NewStorage("authorization", firstAuthorization),
	}

	defaultAuthorizationResponseHandler = authorizationResponseHandler{
		authorizationMgr:       &defaultAuthorizationMgr,
		currentUserFromContext: usermgr.FromContext,
	}
}

func NewHandler() http.Handler {
	route := chi.NewRouter()

	authorizationRequestHandler := authorizationRequestHandler{
		clientByClientID: clientByClientID,
		loginUrl:         url.URL{Path: "/login"},
		authorizationMgr: &defaultAuthorizationMgr,
	}
	route.Get("/auth", authorizationRequestHandler.ServeHTTP)

	route.With(usermgr.Middleware).Get("/callback", defaultAuthorizationResponseHandler.ServeHTTP)

	tokenHandler := &tokenHandler{
		clientByClientID: clientByClientID,
		authorizationMgr: &defaultAuthorizationMgr,
	}
	route.Post("/token", tokenHandler.ServeHTTP)

	return route
}

func RedirectToCallbackWhenLoggedIn(next http.Handler) http.Handler {
	return chi.Chain(defaultAuthorizationMgr.Handler, defaultAuthorizationResponseHandler.redirectWhenLoggedIn).Handler(next)
}
