package server

import (
	"crypto/sha256"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func createSessionToken(u string, t int64, s string) *sessionToken {
	r := sessionToken{UUID: uuid.MustParse(u), CreatedAt: time.Unix(t, 0)}
	r.Salt = [4]byte([]byte(s))
	return &r
}

func TestNewSessionToken(t *testing.T) {
	d := sessionToken{}
	s := newSessionToken()

	if s.CreatedAt == d.CreatedAt {
		t.Errorf("CreatedAt has not be initialized")
	}

	if s.UUID == d.UUID {
		t.Errorf("CreatedAt has not be initialized")
	}

	if s.Salt == d.Salt {
		t.Errorf("CreatedAt has not be initialized")
	}
}

func TestMarshalBinarySessionToken(t *testing.T) {
	orig := newSessionToken()
	b, err := orig.MarshalBinary()
	if err != nil {
		t.Fatalf("failed with %v", err)
	}

	new := &sessionToken{}
	new.UnmarshalBinary(b)

	if !reflect.DeepEqual(orig, new) {
		t.Errorf("orig and new are not the same:\norig=%#v\nnew=%#v", orig, new)
	}
}

func TestSignedSessionToken(t *testing.T) {
	key := []byte("secure!!!")

	for expected, s := range map[string]*sessionToken{
		"AAAAAAAAAAAAAAAAAAAAAAAAAACIbgkAlwWFXLNd7Xfs_ZdTMKkkcnvxrwPKlXHp0qrMTV1EABI": &sessionToken{},
		"YWJjZMtGR_c2iUvMuF29G_7xiFUAAAAAoi54y4EnLqb5ggEn6ng7FKxwXw8-iAGaVdoBnas6R3A": createSessionToken("cb4647f736894bccb85dbd1bfef18855", 0, "abcd"),
		"YWJjZOlyXRxo6EHTiGJ1YbtU35MAAAAAlcvsrJ7mWaLv2pFHwdjIt0JWjhUVVn0xs7dw0XcmD5Q": createSessionToken("e9725d1c68e841d388627561bb54df93", 0, "abcd"),
		"YWJjZP6olW1VlUo3otsekH6RMOoAAAAAOEA0f1P6S41-UHvrY6zb6ZkoSGviYWZDFOcmUUZ5HPY": createSessionToken("fea8956d55954a37a2db1e907e9130ea", 0, "abcd"),
		"YWJjZMtGR_c2iUvMuF29G_7xiFUAAAABh0AaIH5nlL2E6WSvzava009tF6e80vTNNz89gTUQ4lA": createSessionToken("cb4647f736894bccb85dbd1bfef18855", 1, "abcd"),
		"YWJjZOlyXRxo6EHTiGJ1YbtU35MAAAABZeKq2xA2uiHvnrqRmgLZbVNCwM9xq7qX1Qjv9FUKz9c": createSessionToken("e9725d1c68e841d388627561bb54df93", 1, "abcd"),
		"YWJjZP6olW1VlUo3otsekH6RMOoAAAABc-kF8MXX_7kwxuBny9keZVmwiwpKZsil-XWkyKk5c5o": createSessionToken("fea8956d55954a37a2db1e907e9130ea", 1, "abcd"),
		"MTIzNMtGR_c2iUvMuF29G_7xiFUAAAAAS6LGszEa5goQjAEACO4Sx7-XilYa27CxeQuAlW9jhNQ": createSessionToken("cb4647f736894bccb85dbd1bfef18855", 0, "1234"),
		"MTIzNOlyXRxo6EHTiGJ1YbtU35MAAAAAa8aVmKVfnk5rr6mIiGbCpB4J4Xg2QuT7bqQAOLPjqxc": createSessionToken("e9725d1c68e841d388627561bb54df93", 0, "1234"),
		"MTIzNMtGR_c2iUvMuF29G_7xiFUAAAAB3Z2dMSJNyh9fOdH_PAXyDdeOBW5yVskh1sFvYgSJzjI": createSessionToken("cb4647f736894bccb85dbd1bfef18855", 1, "1234"),
		"MTIzNOlyXRxo6EHTiGJ1YbtU35MAAAABBbfrcwqBRRsAO9JdIIxlUgSWlvR0knstidChTyBP0Lg": createSessionToken("e9725d1c68e841d388627561bb54df93", 1, "1234"),
	} {
		if got, _ := s.signedToken(key, sha256.New); got != expected {
			t.Errorf("Expected result %#v not %#v", expected, got)
		}
	}
}

func TestParseSessionToken(t *testing.T) {
	key := []byte("secure!!!")

	for token, expected := range map[string]*sessionToken{
		"YWJjZMtGR_c2iUvMuF29G_7xiFUAAAAAoi54y4EnLqb5ggEn6ng7FKxwXw8-iAGaVdoBnas6R3A": createSessionToken("cb4647f736894bccb85dbd1bfef18855", 0, "abcd"),
		"YWJjZOlyXRxo6EHTiGJ1YbtU35MAAAAAlcvsrJ7mWaLv2pFHwdjIt0JWjhUVVn0xs7dw0XcmD5Q": createSessionToken("e9725d1c68e841d388627561bb54df93", 0, "abcd"),
		"YWJjZP6olW1VlUo3otsekH6RMOoAAAAAOEA0f1P6S41-UHvrY6zb6ZkoSGviYWZDFOcmUUZ5HPY": createSessionToken("fea8956d55954a37a2db1e907e9130ea", 0, "abcd"),
		"YWJjZMtGR_c2iUvMuF29G_7xiFUAAAABh0AaIH5nlL2E6WSvzava009tF6e80vTNNz89gTUQ4lA": createSessionToken("cb4647f736894bccb85dbd1bfef18855", 1, "abcd"),
		"YWJjZOlyXRxo6EHTiGJ1YbtU35MAAAABZeKq2xA2uiHvnrqRmgLZbVNCwM9xq7qX1Qjv9FUKz9c": createSessionToken("e9725d1c68e841d388627561bb54df93", 1, "abcd"),
		"YWJjZP6olW1VlUo3otsekH6RMOoAAAABc-kF8MXX_7kwxuBny9keZVmwiwpKZsil-XWkyKk5c5o": createSessionToken("fea8956d55954a37a2db1e907e9130ea", 1, "abcd"),
		"MTIzNMtGR_c2iUvMuF29G_7xiFUAAAAAS6LGszEa5goQjAEACO4Sx7-XilYa27CxeQuAlW9jhNQ": createSessionToken("cb4647f736894bccb85dbd1bfef18855", 0, "1234"),
		"MTIzNOlyXRxo6EHTiGJ1YbtU35MAAAAAa8aVmKVfnk5rr6mIiGbCpB4J4Xg2QuT7bqQAOLPjqxc": createSessionToken("e9725d1c68e841d388627561bb54df93", 0, "1234"),
		"MTIzNMtGR_c2iUvMuF29G_7xiFUAAAAB3Z2dMSJNyh9fOdH_PAXyDdeOBW5yVskh1sFvYgSJzjI": createSessionToken("cb4647f736894bccb85dbd1bfef18855", 1, "1234"),
		"MTIzNOlyXRxo6EHTiGJ1YbtU35MAAAABBbfrcwqBRRsAO9JdIIxlUgSWlvR0knstidChTyBP0Lg": createSessionToken("e9725d1c68e841d388627561bb54df93", 1, "1234"),
	} {
		got, err := parseToken(key, sha256.New, token)
		if err != nil {
			t.Errorf("parseToken failed: %v", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected result %#v not %#v", expected, got)
		}
	}
}

func FuzzSessionToken(f *testing.F) {
	for i := 0; i < 100; i++ {

		f.Add(rand.Int())
	}

	f.Fuzz(func(t *testing.T, seed int) {
		gen := rand.New(rand.NewSource(int64(seed)))

		salt := make([]byte, 4)
		if _, err := gen.Read(salt); err != nil {
			t.Errorf("failed to setup test: %v", err)
		}

		key := make([]byte, 16)
		if _, err := gen.Read(key); err != nil {
			t.Errorf("failed to setup test: %v", err)
		}

		orig := &sessionToken{
			UUID: uuid.Must(uuid.NewRandomFromReader(gen)), CreatedAt: time.Unix(int64(gen.Uint32()), 0),
			Salt: [4]byte(salt),
		}

		token, err := orig.signedToken(key, sha256.New)
		if err != nil {
			t.Fatalf("failed encoding session %v", err)
		}

		new, err := parseToken(key, sha256.New, token)
		if err != nil {
			t.Fatalf("failed parse token %v", err)
		}

		if !reflect.DeepEqual(orig, new) {
			t.Errorf("sessions are %#v != %#v", orig, new)
		}
	})
}
