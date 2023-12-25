package model

import (
	"math/rand"
	"testing"

	"github.com/seb-schulz/onegate/internal/config"
)

func Test_generateToken(t *testing.T) {
	oldKey := config.Default.Session.Key
	defer func() {
		config.Default.Session.Key = oldKey
	}()
	config.Default.Session.Key = ".test."

	for id, expected := range map[uint]string{1: "1-A-20d74d958f14f6a8bc5f6e7567e32df12b6a5bee55a1bcb037aaeaccc4ce1f51", 2: "2-A-d934bfd86451679cdf3b53f0d885d5454b41859f449cdd94d4049ee803c2934f", 11: "b-A-abf0c4a29c2367331c408441039f72549e5066c259b79d60f125911eb231bfe4", 15: "f-A-9482d52fff73dbe738395c40465a08664dd896c1310533d9cf36f0d5f1518d95", 16: "10-A-472dda9635a9d29e4ebd3a1eff41494f88d53538982014b1e041896229abaebe"} {
		if got := generateToken(id, []byte("A")); got != expected {
			t.Errorf("s.Token = %s; wanted %s", got, expected)
		}
	}
}

func Test_getSessionIDByToken(t *testing.T) {
	oldKey := config.Default.Session.Key
	defer func() {
		config.Default.Session.Key = oldKey
	}()
	config.Default.Session.Key = ".test."

	for expected, token := range map[uint]string{1: "1-A-20d74d958f14f6a8bc5f6e7567e32df12b6a5bee55a1bcb037aaeaccc4ce1f51", 2: "2-A-d934bfd86451679cdf3b53f0d885d5454b41859f449cdd94d4049ee803c2934f", 11: "b-A-abf0c4a29c2367331c408441039f72549e5066c259b79d60f125911eb231bfe4", 15: "f-A-9482d52fff73dbe738395c40465a08664dd896c1310533d9cf36f0d5f1518d95", 16: "10-A-472dda9635a9d29e4ebd3a1eff41494f88d53538982014b1e041896229abaebe"} {
		got, err := getSessionIDByToken(token)
		if err != nil {
			t.Errorf("Got for %v an error: %v", token, err)
		} else if got != expected {
			t.Errorf("s.Token = %x; wanted %x", got, expected)
		}
	}

	for _, token := range []string{"", "1", "x", "1-a", "2-abc-def", "1-B-20d74d958f14f6a8bc5f6e7567e32df12b6a5bee55a1bcb037aaeaccc4ce1f51", "2-A-20d74d958f14f6a8bc5f6e7567e32df12b6a5bee55a1bcb037aaeaccc4ce1f51", "1-A-20d74d958f14f6a8bc5f6e7567e32df12b6a5bee55a1bcb037aaeaccc4ce1f52"} {

		if _, err := getSessionIDByToken(token); err == nil {
			t.Errorf("got no error for %v", token)
		}
	}
}

func FuzzSessionToken(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(uint(rand.Int()))
	}
	f.Fuzz(func(t *testing.T, id uint) {
		s := Session{}
		s.ID = id
		token := s.Token()
		out, err := getSessionIDByToken(token)
		if err != nil {
			t.Fatalf("%v: decode: %v", id, err)
		}
		if id != out {
			t.Fatalf("%v: not equal after round trip: %v", id, out)
		}
	})
}
