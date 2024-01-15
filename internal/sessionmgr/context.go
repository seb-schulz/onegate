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
