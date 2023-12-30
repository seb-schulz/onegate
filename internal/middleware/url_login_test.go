package middleware

import (
	"math/rand"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/seb-schulz/onegate/internal/config"
)

func parseUrl(url *url.URL) error {
	pathChunks := strings.Split(url.Path, "/")
	_, err := parseToken(pathChunks[len(pathChunks)-1])
	return err
}

func FuzzGetLoginUrl(f *testing.F) {
	oldKey := config.Default.UrlLogin.Key
	defer func() {
		config.Default.UrlLogin.Key = oldKey
	}()
	config.Default.UrlLogin.Key = []byte(".test.")

	for i := 0; i < 100; i++ {
		f.Add(uint(rand.Int()))
	}
	f.Fuzz(func(t *testing.T, id uint) {
		out, err := GetLoginUrl(id, time.Second)
		if err != nil {
			t.Fatalf("URL for %v: failed: %v", id, err)
		}

		if err := parseUrl(out); err != nil {
			t.Fatalf("cannot parse URL %v: %v\n", *out, err)
		}
	})
}
