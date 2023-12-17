package graph

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/model"
)

func mustRandomEncodedBytes(len int) string {
	r := make([]byte, len)

	_, err := rand.Read(r)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(r)
}

func mustSessionFromContext(ctx context.Context) *model.Session {
	session := middleware.SessionFromContext(ctx)
	if session == nil {
		panic("session is missing")
	}
	return session
}
