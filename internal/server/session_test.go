package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func createSessionToken(u string, t int64, s string) sessionTokenizer {
	r := sessionToken{UUID: uuid.MustParse(u), CreatedAt: time.Unix(t, 0), sig: sha256.New}
	r.Salt = [4]byte([]byte(s))
	return &r
}

func TestNewSessionToken(t *testing.T) {
	d := sessionToken{}
	s := sessionToken{}
	s.initialize()

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
	orig := sessionToken{}
	orig.initialize()

	b, err := orig.MarshalBinary()
	if err != nil {
		t.Fatalf("failed with %v", err)
	}

	new := sessionToken{}
	new.UnmarshalBinary(b)

	if !reflect.DeepEqual(orig, new) {
		t.Errorf("orig and new are not the same:\norig=%#v\nnew=%#v", orig, new)
	}
}

func TestSignedSessionToken(t *testing.T) {
	key := []byte("secure!!!")

	for expected, s := range map[string]sessionTokenizer{
		"AAAAAAAAAAAAAAAAAAAAAAAAAACIbgkAlwWFXLNd7Xfs_ZdTMKkkcnvxrwPKlXHp0qrMTV1EABI": newSessionToken(),
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
		if got, _ := s.sign(key); expected != base64.RawURLEncoding.EncodeToString(got) {
			t.Errorf("Expected result %#v not %#v", expected, got)
		}
	}
}

func deepEqualSessionToken(a, b *sessionToken) bool {
	return reflect.DeepEqual(a.UUID, b.UUID) && a.CreatedAt.Truncate(time.Second) == b.CreatedAt.Truncate(time.Second) && bytes.Equal(a.Salt[:], b.Salt[:])
}

func TestParseSessionToken(t *testing.T) {
	key := []byte("secure!!!")

	for token, expected := range map[string]sessionTokenizer{
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
		got := newSessionToken()
		err := got.parse(key, rawToken)
		if err != nil {
			t.Errorf("parseToken failed: %v", err)
		}

		if !deepEqualSessionToken(got.(*sessionToken), expected.(*sessionToken)) {
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
			sig:  sha256.New,
		}

		token, err := orig.sign(key)
		if err != nil {
			t.Fatalf("failed encoding session %v", err)
		}

		new := newSessionToken()
		if err := new.parse(key, token); err != nil {
			t.Fatalf("failed parse token %v", err)
		}

		if !deepEqualSessionToken(orig, new.(*sessionToken)) {
			t.Errorf("sessions are %#v != %#v", orig, new)
		}
	})
}

type mockSession struct {
	data    []byte
	signFn  func(s *mockSession, key []byte) ([]byte, error)
	parseFn func(s *mockSession, key []byte, token []byte) error
	initFn  func(*mockSession)
}

func (s *mockSession) parse(key []byte, token []byte) error {
	return s.parseFn(s, key, token)
}

func (s *mockSession) sign(key []byte) ([]byte, error) {
	return s.signFn(s, key)
}

func (s *mockSession) initialize() {
	s.initFn(s)
}

func (s *mockSession) String() string {
	return string(s.data)
}

func foreachCookie(resp *http.Response, name string, fn func(*http.Cookie)) {
	for _, c := range resp.Cookies() {
		if c.Name == "session" {
			fn(c)
		}
	}
}

func newCustomRequest(fns ...func(*http.Request)) *http.Request {
	req := httptest.NewRequest("GET", "/foo", nil)
	for _, fn := range fns {
		fn(req)
	}
	return req
}

func FuzzSessionMiddleware_responds_code(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(rand.Int())
	}

	f.Fuzz(func(t *testing.T, seed int) {

		t.Run("new token must be created", func(t *testing.T) {
			seed := seed
			t.Log(seed)
			t.Parallel()
			gen := rand.New(rand.NewSource(int64(seed)))
			var (
				key   [4]byte
				token [3]byte
			)

			gen.Read(key[:])

			middleware := sessionMiddleware{
				key: key[:],
				newToken: func() sessionTokenizer {
					return &mockSession{
						signFn: func(s *mockSession, k []byte) ([]byte, error) {
							if !bytes.Equal(key[:], k) {
								t.Errorf("keys are not equal %#v != %#v", key[:], k)
							}
							return s.data, nil
						},
						parseFn: func(s *mockSession, key []byte, token []byte) error {
							s.data = token
							return nil
						},
						initFn: func(s *mockSession) {
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Ok")
			}))

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, newCustomRequest())

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.FailNow()
			}

			counter := 0
			foreachCookie(resp, "session", func(cookie *http.Cookie) {
				counter++
				if got, err := base64.RawURLEncoding.DecodeString(cookie.Value); !bytes.Equal(got, token[:]) {
					t.Errorf("expected %#v but got %v (%#v) with %v", token[:], got, cookie.Value, err)
				}
			})

			if counter != 1 {
				t.Errorf("countend %d cookies instead of 1", counter)
			}
		})

		t.Run("existing invalid token", func(t *testing.T) {
			seed := seed
			t.Log(seed)

			t.Parallel()

			gen := rand.New(rand.NewSource(int64(seed)))
			var (
				key          [4]byte
				token        [3]byte
				invalidToken [12]byte
			)

			gen.Read(key[:])

			middleware := sessionMiddleware{
				key: key[:],
				newToken: func() sessionTokenizer {
					return &mockSession{
						signFn: func(s *mockSession, k []byte) ([]byte, error) {
							if !bytes.Equal(key[:], k) {
								t.Errorf("keys are not equal %#v != %#v", key[:], k)
							}
							return s.data, nil
						},
						parseFn: func(s *mockSession, key []byte, token []byte) error {
							return fmt.Errorf("invalid token")
						},
						initFn: func(s *mockSession) {
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Ok")
			}))

			gen.Read(invalidToken[:])

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, newCustomRequest(func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:     "session",
					Value:    string(invalidToken[:]),
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
			}))

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.FailNow()
			}

			counter := 0
			foreachCookie(resp, "session", func(cookie *http.Cookie) {
				counter++
				if got, err := base64.RawURLEncoding.DecodeString(cookie.Value); !bytes.Equal(got, token[:]) {
					t.Errorf("expected %#v but got %v (%#v) with %v", token[:], got, cookie.Value, err)
				}
			})

			if counter != 1 {
				t.Errorf("countend %d cookies instead of 1", counter)
			}
		})

		t.Run("existing valid token", func(t *testing.T) {
			seed := seed
			t.Log(seed)

			t.Parallel()

			gen := rand.New(rand.NewSource(int64(seed)))
			var (
				key   [4]byte
				token [3]byte
			)

			gen.Read(key[:])

			middleware := sessionMiddleware{
				key: key[:],
				newToken: func() sessionTokenizer {
					return &mockSession{
						signFn: func(s *mockSession, k []byte) ([]byte, error) {
							t.Error("this func does not need to be called")
							return s.data, nil
						},
						parseFn: func(s *mockSession, key []byte, got []byte) error {
							if !bytes.Equal(token[:], got) {
								t.Errorf("Got %#v instead of %#v.", got, token[:])
							}
							s.data = got
							return nil
						},
						initFn: func(s *mockSession) {
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Ok")
			}))

			gen.Read(token[:])

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, newCustomRequest(func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:     "session",
					Value:    base64.RawURLEncoding.EncodeToString(token[:]),
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
			}))

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.FailNow()
			}

			counter := 0
			foreachCookie(resp, "session", func(cookie *http.Cookie) {
				counter++
			})

			if counter > 0 {
				t.Errorf("countend %d cookies instead of zero", counter)
			}
		})
	})
}
func FuzzSessionMiddleware_context(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(rand.Int())
	}

	f.Fuzz(func(t *testing.T, seed int) {

		t.Run("new token must be created", func(t *testing.T) {
			seed := seed
			t.Parallel()
			gen := rand.New(rand.NewSource(int64(seed)))
			var (
				key   [4]byte
				token [3]byte
			)

			gen.Read(key[:])

			middleware := sessionMiddleware{
				key: key[:],
				newToken: func() sessionTokenizer {
					return &mockSession{
						signFn: func(s *mockSession, k []byte) ([]byte, error) {
							if !bytes.Equal(key[:], k) {
								t.Errorf("keys are not equal %#v != %#v", key[:], k)
							}
							return s.data, nil
						},
						parseFn: func(s *mockSession, key []byte, token []byte) error {
							s.data = token
							return nil
						},
						initFn: func(s *mockSession) {
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionToken := mustSessionTokenFromContext(r.Context())
				if !bytes.Equal(sessionToken.(*mockSession).data, token[:]) {
					t.Errorf("Got %#v instead of %#v", sessionToken.(*mockSession).data, token[:])
				}
				fmt.Fprintln(w, "Ok")
			}))

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, newCustomRequest())
		})

		t.Run("existing invalid token", func(t *testing.T) {
			seed := seed
			t.Log(seed)

			t.Parallel()

			gen := rand.New(rand.NewSource(int64(seed)))
			var (
				key          [4]byte
				token        [3]byte
				invalidToken [12]byte
			)

			gen.Read(key[:])

			sessionMiddleware := sessionMiddleware{
				key: key[:],
				newToken: func() sessionTokenizer {
					return &mockSession{
						signFn: func(s *mockSession, k []byte) ([]byte, error) {
							if !bytes.Equal(key[:], k) {
								t.Errorf("keys are not equal %#v != %#v", key[:], k)
							}
							return s.data, nil
						},
						parseFn: func(s *mockSession, key []byte, token []byte) error {
							return fmt.Errorf("invalid token")
						},
						initFn: func(s *mockSession) {
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := sessionMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionToken := mustSessionTokenFromContext(r.Context())
				if !bytes.Equal(sessionToken.(*mockSession).data, token[:]) {
					t.Errorf("Got %#v instead of %#v", sessionToken.(*mockSession).data, token[:])
				}
				fmt.Fprintln(w, "Ok")
			}))
			// handler = middleware.Logger(handler)

			gen.Read(invalidToken[:])

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, newCustomRequest(func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:     "session",
					Value:    string(invalidToken[:]),
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
			}))
		})

		t.Run("existing valid token", func(t *testing.T) {
			seed := seed
			t.Log(seed)

			t.Parallel()

			gen := rand.New(rand.NewSource(int64(seed)))
			var (
				key   [4]byte
				token [3]byte
			)

			gen.Read(key[:])

			middleware := sessionMiddleware{
				key: key[:],
				newToken: func() sessionTokenizer {
					return &mockSession{
						signFn: func(s *mockSession, k []byte) ([]byte, error) {
							t.Error("this func does not need to be called")
							return s.data, nil
						},
						parseFn: func(s *mockSession, key []byte, got []byte) error {
							if !bytes.Equal(token[:], got) {
								t.Errorf("Got %#v instead of %#v.", got, token[:])
							}
							s.data = got
							return nil
						},
						initFn: func(s *mockSession) {
							t.Error("this func does not need to be called")
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionToken := mustSessionTokenFromContext(r.Context())
				if !bytes.Equal(sessionToken.(*mockSession).data, token[:]) {
					t.Errorf("Got %#v instead of %#v", sessionToken.(*mockSession).data, token[:])
				}
				fmt.Fprintln(w, "Ok")
			}))

			gen.Read(token[:])

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, newCustomRequest(func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:     "session",
					Value:    base64.RawURLEncoding.EncodeToString(token[:]),
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
			}))
		})
	})
}
