package server

import (
	"bytes"
	"context"
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"math/rand"
	"net/http"
	"net/http/httptest"
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
	s := newSessionToken().(*sessionToken)

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
	orig := newSessionToken().(*sessionToken)
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
		if got, _ := s.signedToken(key, sha256.New); expected != base64.RawURLEncoding.EncodeToString(got) {
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
		rawToken, _ := base64.RawURLEncoding.DecodeString(token)
		got, err := parseToken(key, sha256.New, rawToken)
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

type sessionTokenParser func(key []byte, sig func() hash.Hash, token []byte) (sessionTokenizer, error)
type sessionTokenGenerator func() sessionTokenizer
type contextSessionKeyType struct{ string }

var contextSessionToken = contextSessionKeyType{"session"}

type sessionMiddleware struct {
	key            []byte
	tokenGenerator sessionTokenGenerator
	tokenParser    sessionTokenParser
	signer         func() hash.Hash
}

func (s *sessionMiddleware) setCookie(w http.ResponseWriter, token sessionTokenizer) {
	sToken, err := token.signedToken(s.key, s.signer)
	if err != nil {
		panic(err) // signing token should not fail
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    base64.RawURLEncoding.EncodeToString(sToken),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (s *sessionMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token sessionTokenizer

		cookie, err := r.Cookie("session")
		if err != nil {
			token = s.tokenGenerator()
			s.setCookie(w, token)
			ctx := context.WithValue(r.Context(), contextSessionToken, token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		rawToken, err := base64.RawURLEncoding.DecodeString(cookie.Value)
		if err != nil {
			token = s.tokenGenerator()
			s.setCookie(w, token)
			// ctx := context.WithValue(r.Context(), contextSessionToken, token)
			// next.ServeHTTP(w, r.WithContext(ctx))
			next.ServeHTTP(w, r)
			return
		}

		token, err = s.tokenParser(s.key, s.signer, rawToken)
		if err != nil {
			token = s.tokenGenerator()
			s.setCookie(w, token)
			ctx := context.WithValue(r.Context(), contextSessionToken, token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), contextSessionToken, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func mustSessionTokenFromContext(ctx context.Context) sessionTokenizer {
	raw, ok := ctx.Value(contextSessionToken).(sessionTokenizer)
	if !ok {
		panic("session token does not exist")
	}
	return raw
}

type mockSession struct {
	data []byte
}

func parseMockedToken(key []byte, sig func() hash.Hash, token []byte) (sessionTokenizer, error) {
	h := hmac.New(sig, key)
	idx := len(token) - sig().Size()
	if idx <= 0 {
		return nil, fmt.Errorf("invalid")
	}
	payload, signature := token[:idx], token[idx:]

	if _, err := h.Write(payload); err != nil {
		return nil, err
	}

	if !hmac.Equal(signature, h.Sum(nil)) {
		return nil, fmt.Errorf("token is tampered")
	}

	return &mockSession{data: payload}, nil
}

func (s *mockSession) signedToken(key []byte, sig func() hash.Hash) ([]byte, error) {
	h := hmac.New(sig, key)

	if _, err := h.Write(s.data); err != nil {
		return []byte{}, err
	}

	return append(s.data, h.Sum(nil)...), nil
}

func foreachCookie(resp *http.Response, name string, fn func(*http.Cookie)) {
	for _, c := range resp.Cookies() {
		if c.Name == "session" {
			fn(c)
		}
	}
}

func newCustomRequest(fns ...func(*http.Request)) *http.Request {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	for _, fn := range fns {
		fn(req)
	}
	return req
}

func randToken() []byte {
	var r [4]byte
	if _, err := crand.Read(r[:]); err != nil {
		panic(err)
	}

	return r[:]
}

func TestSessionMiddleware(t *testing.T) {
	type scenario struct {
		token        []byte
		req          *http.Request
		handlerCheck http.HandlerFunc
		respCheck    func(*http.Response)
	}

	key := []byte(".test.")
	for _, testCase := range []scenario{
		{
			[]byte("abcd"),
			newCustomRequest(),
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Ok")
				token := mustSessionTokenFromContext(r.Context())
				if !bytes.Equal(token.(*mockSession).data, []byte("abcd")) {
					t.Errorf("token are not equal: %#v", token)
				}
			},
			func(resp *http.Response) {
				if resp.StatusCode != http.StatusOK {
					t.FailNow()
				}

				counter := 0
				foreachCookie(resp, "session", func(cookie *http.Cookie) {
					counter++
					if cookie.Value != "YWJjZP0k6SMGsr8u3svFp-yHN-LRwBbr" {
						t.Errorf("cookie value invalid: %v", cookie.Value)
					}
				})

				if counter != 1 {
					t.Errorf("countend %d cookies instead of 1", counter)
				}
			},
		}, {
			[]byte("efgh"),
			newCustomRequest(func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:     "session",
					Value:    "invalidToken",
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
			}),
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Ok")

				token := mustSessionTokenFromContext(r.Context())
				if !bytes.Equal(token.(*mockSession).data, []byte("efgh")) {
					t.Errorf("token are not equal: %#v", token)
				}
			},
			func(resp *http.Response) {
				if resp.StatusCode != http.StatusOK {
					t.FailNow()
				}

				counter := 0
				foreachCookie(resp, "session", func(cookie *http.Cookie) {
					counter++
					if cookie.Value != "ZWZnaEo_hvbODdyh623e7aVXM1NBh-yS" {
						t.Errorf("cookie value invalid: %v", cookie.Value)
					}
				})

				if counter != 1 {
					t.Errorf("countend %d cookies instead of 1", counter)
				}
			},
		}, {
			randToken(),
			newCustomRequest(func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:     "session",
					Value:    "YWJjJONbU5tdpIfhErhiyEMCELIc47U",
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
			}),
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Ok")

				token := mustSessionTokenFromContext(r.Context())
				if !bytes.Equal(token.(*mockSession).data, []byte("abc")) {
					t.Errorf("token are not equal: %#v", token)
				}
			},
			func(resp *http.Response) {
				if resp.StatusCode != http.StatusOK {
					t.FailNow()
				}
				counter := 0
				foreachCookie(resp, "session", func(cookie *http.Cookie) {
					counter++
				})

				if counter > 0 {
					t.Errorf("countend %d cookies instead of 0", counter)
				}
			},
		},
	} {
		newMockSession := func() sessionTokenizer {
			return &mockSession{testCase.token}
		}

		middleware := sessionMiddleware{
			key:            key,
			tokenGenerator: newMockSession,
			tokenParser:    parseMockedToken,
			signer:         sha1.New,
		}

		handler := middleware.Handler(http.HandlerFunc(testCase.handlerCheck))

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, testCase.req)

		resp := w.Result()
		testCase.respCheck(resp)
	}
}
