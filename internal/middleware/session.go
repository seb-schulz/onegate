package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/seb-schulz/onegate/internal/model"
	"gorm.io/gorm"
)

type contextSessionKeyType int

const (
	contextSessionKey contextSessionKeyType = iota
)

func SessionMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var session model.Session
			tx := db.WithContext(context.Background()).Debug()
			cookie, err := r.Cookie("session")
			if err != nil {
				model.CreateSession(tx, &session)
				http.SetCookie(w, &http.Cookie{
					Name:    "session",
					Value:   session.Token(),
					Expires: time.Now().Add(30 * time.Minute),
				})
			} else {
				err := model.FirstSessionByToken(tx, cookie.Value, &session)
				if err != nil {
					panic(err)
				}

			}

			ctx := context.WithValue(r.Context(), contextSessionKey, &session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SessionFromContext(ctx context.Context) *model.Session {
	raw, _ := ctx.Value(contextSessionKey).(*model.Session)
	return raw
}
