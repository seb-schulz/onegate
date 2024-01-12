package sessionmgr

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func createToken(u string, t int64, s string) tokenizer {
	r := Token{UUID: uuid.MustParse(u), CreatedAt: time.Unix(t, 0), sig: sha256.New}
	r.salt = [4]byte([]byte(s))
	return &r
}

func TestNewToken(t *testing.T) {
	d := Token{}
	s := Token{}
	s.initialize()

	if s.CreatedAt == d.CreatedAt {
		t.Errorf("CreatedAt has not be initialized")
	}

	if s.UUID == d.UUID {
		t.Errorf("CreatedAt has not be initialized")
	}

	if s.salt == d.salt {
		t.Errorf("CreatedAt has not be initialized")
	}
}

func TestMarshalBinaryToken(t *testing.T) {
	orig := Token{}
	orig.initialize()

	b, err := orig.MarshalBinary()
	if err != nil {
		t.Fatalf("failed with %v", err)
	}

	new := Token{}
	new.UnmarshalBinary(b)

	if !reflect.DeepEqual(orig, new) {
		t.Errorf("orig and new are not the same:\norig=%#v\nnew=%#v", orig, new)
	}
}

func TestSignedToken(t *testing.T) {
	key := []byte("secure!!!")

	for expected, s := range map[string]tokenizer{
		"AAAAAAAAAAAAAAAAAAAAAAAAAACIbgkAlwWFXLNd7Xfs_ZdTMKkkcnvxrwPKlXHp0qrMTV1EABI": newToken(),
		"YWJjZMtGR_c2iUvMuF29G_7xiFUAAAAAoi54y4EnLqb5ggEn6ng7FKxwXw8-iAGaVdoBnas6R3A": createToken("cb4647f736894bccb85dbd1bfef18855", 0, "abcd"),
		"YWJjZOlyXRxo6EHTiGJ1YbtU35MAAAAAlcvsrJ7mWaLv2pFHwdjIt0JWjhUVVn0xs7dw0XcmD5Q": createToken("e9725d1c68e841d388627561bb54df93", 0, "abcd"),
		"YWJjZP6olW1VlUo3otsekH6RMOoAAAAAOEA0f1P6S41-UHvrY6zb6ZkoSGviYWZDFOcmUUZ5HPY": createToken("fea8956d55954a37a2db1e907e9130ea", 0, "abcd"),
		"YWJjZMtGR_c2iUvMuF29G_7xiFUAAAABh0AaIH5nlL2E6WSvzava009tF6e80vTNNz89gTUQ4lA": createToken("cb4647f736894bccb85dbd1bfef18855", 1, "abcd"),
		"YWJjZOlyXRxo6EHTiGJ1YbtU35MAAAABZeKq2xA2uiHvnrqRmgLZbVNCwM9xq7qX1Qjv9FUKz9c": createToken("e9725d1c68e841d388627561bb54df93", 1, "abcd"),
		"YWJjZP6olW1VlUo3otsekH6RMOoAAAABc-kF8MXX_7kwxuBny9keZVmwiwpKZsil-XWkyKk5c5o": createToken("fea8956d55954a37a2db1e907e9130ea", 1, "abcd"),
		"MTIzNMtGR_c2iUvMuF29G_7xiFUAAAAAS6LGszEa5goQjAEACO4Sx7-XilYa27CxeQuAlW9jhNQ": createToken("cb4647f736894bccb85dbd1bfef18855", 0, "1234"),
		"MTIzNOlyXRxo6EHTiGJ1YbtU35MAAAAAa8aVmKVfnk5rr6mIiGbCpB4J4Xg2QuT7bqQAOLPjqxc": createToken("e9725d1c68e841d388627561bb54df93", 0, "1234"),
		"MTIzNMtGR_c2iUvMuF29G_7xiFUAAAAB3Z2dMSJNyh9fOdH_PAXyDdeOBW5yVskh1sFvYgSJzjI": createToken("cb4647f736894bccb85dbd1bfef18855", 1, "1234"),
		"MTIzNOlyXRxo6EHTiGJ1YbtU35MAAAABBbfrcwqBRRsAO9JdIIxlUgSWlvR0knstidChTyBP0Lg": createToken("e9725d1c68e841d388627561bb54df93", 1, "1234"),
	} {
		if got, _ := s.sign(key); expected != base64.RawURLEncoding.EncodeToString(got) {
			t.Errorf("Expected result %#v not %#v", expected, got)
		}
	}
}

func deepEqualToken(a, b *Token) bool {
	return reflect.DeepEqual(a.UUID, b.UUID) && a.CreatedAt.Truncate(time.Second) == b.CreatedAt.Truncate(time.Second) && bytes.Equal(a.salt[:], b.salt[:])
}

func TestParseToken(t *testing.T) {
	key := []byte("secure!!!")

	for token, expected := range map[string]tokenizer{
		"YWJjZMtGR_c2iUvMuF29G_7xiFUAAAAAoi54y4EnLqb5ggEn6ng7FKxwXw8-iAGaVdoBnas6R3A": createToken("cb4647f736894bccb85dbd1bfef18855", 0, "abcd"),
		"YWJjZOlyXRxo6EHTiGJ1YbtU35MAAAAAlcvsrJ7mWaLv2pFHwdjIt0JWjhUVVn0xs7dw0XcmD5Q": createToken("e9725d1c68e841d388627561bb54df93", 0, "abcd"),
		"YWJjZP6olW1VlUo3otsekH6RMOoAAAAAOEA0f1P6S41-UHvrY6zb6ZkoSGviYWZDFOcmUUZ5HPY": createToken("fea8956d55954a37a2db1e907e9130ea", 0, "abcd"),
		"YWJjZMtGR_c2iUvMuF29G_7xiFUAAAABh0AaIH5nlL2E6WSvzava009tF6e80vTNNz89gTUQ4lA": createToken("cb4647f736894bccb85dbd1bfef18855", 1, "abcd"),
		"YWJjZOlyXRxo6EHTiGJ1YbtU35MAAAABZeKq2xA2uiHvnrqRmgLZbVNCwM9xq7qX1Qjv9FUKz9c": createToken("e9725d1c68e841d388627561bb54df93", 1, "abcd"),
		"YWJjZP6olW1VlUo3otsekH6RMOoAAAABc-kF8MXX_7kwxuBny9keZVmwiwpKZsil-XWkyKk5c5o": createToken("fea8956d55954a37a2db1e907e9130ea", 1, "abcd"),
		"MTIzNMtGR_c2iUvMuF29G_7xiFUAAAAAS6LGszEa5goQjAEACO4Sx7-XilYa27CxeQuAlW9jhNQ": createToken("cb4647f736894bccb85dbd1bfef18855", 0, "1234"),
		"MTIzNOlyXRxo6EHTiGJ1YbtU35MAAAAAa8aVmKVfnk5rr6mIiGbCpB4J4Xg2QuT7bqQAOLPjqxc": createToken("e9725d1c68e841d388627561bb54df93", 0, "1234"),
		"MTIzNMtGR_c2iUvMuF29G_7xiFUAAAAB3Z2dMSJNyh9fOdH_PAXyDdeOBW5yVskh1sFvYgSJzjI": createToken("cb4647f736894bccb85dbd1bfef18855", 1, "1234"),
		"MTIzNOlyXRxo6EHTiGJ1YbtU35MAAAABBbfrcwqBRRsAO9JdIIxlUgSWlvR0knstidChTyBP0Lg": createToken("e9725d1c68e841d388627561bb54df93", 1, "1234"),
	} {
		rawToken, _ := base64.RawURLEncoding.DecodeString(token)
		got := newToken()
		err := got.parse(key, rawToken)
		if err != nil {
			t.Errorf("parseToken failed: %v", err)
		}

		if !deepEqualToken(got.(*Token), expected.(*Token)) {
			t.Errorf("Expected result %#v not %#v", expected, got)
		}
	}
}

func FuzzToken(f *testing.F) {
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

		orig := &Token{
			UUID: uuid.Must(uuid.NewRandomFromReader(gen)), CreatedAt: time.Unix(int64(gen.Uint32()), 0),
			salt: [4]byte(salt),
			sig:  sha256.New,
		}

		token, err := orig.sign(key)
		if err != nil {
			t.Fatalf("failed encoding session %v", err)
		}

		new := newToken()
		if err := new.parse(key, token); err != nil {
			t.Fatalf("failed parse token %v", err)
		}

		if !deepEqualToken(orig, new.(*Token)) {
			t.Errorf("sessions are %#v != %#v", orig, new)
		}
	})
}
