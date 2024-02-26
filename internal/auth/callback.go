package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/seb-schulz/onegate/internal/model"
)

type callbackRedirectHandler struct {
	currentAuthorization func(context.Context) (authorization, error)
	currentUser          func(ctx context.Context) *model.User
}

func (cr callbackRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := cr.currentUser(r.Context())
	if user == nil {
		slog.Warn("oAuth2 callback failed because user not logged in")
		http.NotFound(w, r)
		return
	}

	authReq, err := cr.currentAuthorization(r.Context())
	if err != nil {
		slog.Warn(fmt.Sprintf("oAuth2 callback failed with: %v", err))
		http.NotFound(w, r)
		return
	}

	if err := authReq.SetUserID(r.Context(), user.ID); err != nil {
		slog.Error(fmt.Sprintf("Failed to assign user ID: %v", err))
		http.NotFound(w, r)
		return
	}

	q := url.Values{}
	q.Add("code", authReq.Code())
	q.Add("state", authReq.State())
	http.Redirect(w, r, fmt.Sprintf("%v?%v", authReq.RedirectURI(), q.Encode()), http.StatusFound)
}
