package sessionmgr

import (
	"context"
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

func ToContext(ctx context.Context, t *Token) context.Context {
	return context.WithValue(ctx, contextToken, t)
}
