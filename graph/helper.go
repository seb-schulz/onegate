package graph

import (
	"crypto/rand"
	"encoding/base64"
)

func mustRandomEncodedBytes(len int) string {
	r := make([]byte, len)

	_, err := rand.Read(r)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(r)
}
