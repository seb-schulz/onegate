package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/seb-schulz/onegate/internal/model"
)

type authorizationResponseHandler struct {
	currentUserFromContext func(ctx context.Context) *model.User
	authorizationMgr       interface {
		updateUserID(ctx context.Context, userID uint) error
		fromContext(ctx context.Context) authorization
	}
}

func (auth authorizationResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.currentUserFromContext(r.Context())
	if err := auth.authorizationMgr.updateUserID(r.Context(), currentUser.ID); err != nil {
		// An error should not happen because user must be authenticated.
		// Reasoning could be tampered authentication flow or external error like database outage.
		panic(fmt.Errorf("failed to update user id: %w", err))
	}
	authReq := auth.authorizationMgr.fromContext(r.Context())
	q := url.Values{}
	q.Add("code", authReq.Code())
	q.Add("state", authReq.State())
	http.Redirect(w, r, fmt.Sprintf("%v?%v", authReq.RedirectURI(), q.Encode()), http.StatusFound)
}

func (auth authorizationResponseHandler) redirectWhenLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUser := auth.currentUserFromContext(r.Context())
		if currentUser == nil {
			next.ServeHTTP(w, r)
			return
		}

		authReq := auth.authorizationMgr.fromContext(r.Context())
		slog.Info(fmt.Sprintf("Redirect middleware: Redirect: %#v authReq=%#v", authReq.RedirectURI(), authReq))
		if authReq.Exists() {
			auth.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
