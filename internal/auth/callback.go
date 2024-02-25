package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type callbackRedirectHandler struct {
	currentAuthorization func(context.Context) (authorization, error)
}

func (cr callbackRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authReq, err := cr.currentAuthorization(r.Context())
	if err != nil {
		slog.Warn(fmt.Sprintf("oAuth2 callback failed with: %v", err))
		http.NotFound(w, r)
		return
	}

	q := url.Values{}
	q.Add("code", authReq.Code())
	q.Add("state", authReq.State())
	http.Redirect(w, r, fmt.Sprintf("%v?%v", authReq.RedirectURI(), q.Encode()), http.StatusFound)
}
