package sessionmgr

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

type (
	middleware struct {
		key      []byte
		newToken func() tokenizer
	}
)

func (s *middleware) setCookie(w http.ResponseWriter, token tokenSigner) {
	sToken, err := token.sign(s.key)
	if err != nil {
		panic(err) // signing token should not fail
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    base64.RawURLEncoding.EncodeToString(sToken),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (s *middleware) tokenFromCookie(req *http.Request, token tokenParser) error {
	cookie, err := req.Cookie("session")
	if err != nil {
		return err
	}

	rawToken, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return err
	}

	if err := token.parse(s.key, rawToken); err != nil {
		logger := httplog.LogEntry(req.Context())
		logger.Warn(fmt.Sprintf("cannot parse raw token: %v", err))
		return err
	}
	return nil
}

func (s *middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := s.newToken()

		if err := s.tokenFromCookie(r, token); err != nil {
			token.initialize()
			s.setCookie(w, token)
		}

		ctx := context.WithValue(r.Context(), contextToken, token)
		httplog.LogEntrySetField(ctx, "session", slog.StringValue(fmt.Sprint(token)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func DefaultMiddleware(key []byte) func(next http.Handler) http.Handler {
	s := middleware{key, newToken}
	return s.Handler
}
