package sessionmgr

import (
	"context"

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

func ContextWithTransaction[T any](ctx context.Context, tx *gorm.DB, fn func(*gorm.DB, *Token) (T, error)) (T, error) {
	var result T
	err := tx.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = fn(tx, FromContext(ctx))
		return err
	})
	return result, err
}
