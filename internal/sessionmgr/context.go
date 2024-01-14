package sessionmgr

import (
	"context"

	"github.com/seb-schulz/onegate/internal/database"
	"gorm.io/gorm"
)

type contextSessionKeyType struct{ string }

var contextToken = contextSessionKeyType{"session"}

func FromContext(ctx context.Context) *Token {
	raw, ok := ctx.Value(contextToken).(*Token)
	if !ok {
		panic("session token does not exist")
	}
	return raw
}

func ContextWithToken[T any](ctx context.Context, fn func(*gorm.DB, *Token) (T, error)) (T, error) {
	return database.Transaction(ctx, func(tx *gorm.DB) (T, error) {
		return fn(tx, FromContext(ctx))
	})
}
