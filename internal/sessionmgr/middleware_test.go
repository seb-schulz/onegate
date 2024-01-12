package sessionmgr

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockToken struct {
	data    []byte
	signFn  func(s *mockToken, key []byte) ([]byte, error)
	parseFn func(s *mockToken, key []byte, token []byte) error
	initFn  func(*mockToken)
}

func (s *mockToken) parse(key []byte, token []byte) error {
	return s.parseFn(s, key, token)
}

func (s *mockToken) sign(key []byte) ([]byte, error) {
	return s.signFn(s, key)
}

func (s *mockToken) initialize() {
	s.initFn(s)
}

func (s *mockToken) String() string {
	return string(s.data)
}

func foreachCookie(resp *http.Response, name string, fn func(*http.Cookie)) {
	for _, c := range resp.Cookies() {
		if c.Name == "session" {
			fn(c)
		}
	}
}

func fromContext(ctx context.Context) *mockToken {
	raw, ok := ctx.Value(contextToken).(*mockToken)
	if !ok {
		panic("session token does not exist")
	}
	return raw
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

			middleware := middleware{
				key: key[:],
				newToken: func() tokenizer {
					return &mockToken{
						signFn: func(s *mockToken, k []byte) ([]byte, error) {
							if !bytes.Equal(key[:], k) {
								t.Errorf("keys are not equal %#v != %#v", key[:], k)
							}
							return s.data, nil
						},
						parseFn: func(s *mockToken, key []byte, token []byte) error {
							s.data = token
							return nil
						},
						initFn: func(s *mockToken) {
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

			middleware := middleware{
				key: key[:],
				newToken: func() tokenizer {
					return &mockToken{
						signFn: func(s *mockToken, k []byte) ([]byte, error) {
							if !bytes.Equal(key[:], k) {
								t.Errorf("keys are not equal %#v != %#v", key[:], k)
							}
							return s.data, nil
						},
						parseFn: func(s *mockToken, key []byte, token []byte) error {
							return fmt.Errorf("invalid token")
						},
						initFn: func(s *mockToken) {
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

			middleware := middleware{
				key: key[:],
				newToken: func() tokenizer {
					return &mockToken{
						signFn: func(s *mockToken, k []byte) ([]byte, error) {
							t.Error("this func does not need to be called")
							return s.data, nil
						},
						parseFn: func(s *mockToken, key []byte, got []byte) error {
							if !bytes.Equal(token[:], got) {
								t.Errorf("Got %#v instead of %#v.", got, token[:])
							}
							s.data = got
							return nil
						},
						initFn: func(s *mockToken) {
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

			middleware := middleware{
				key: key[:],
				newToken: func() tokenizer {
					return &mockToken{
						signFn: func(s *mockToken, k []byte) ([]byte, error) {
							if !bytes.Equal(key[:], k) {
								t.Errorf("keys are not equal %#v != %#v", key[:], k)
							}
							return s.data, nil
						},
						parseFn: func(s *mockToken, key []byte, token []byte) error {
							s.data = token
							return nil
						},
						initFn: func(s *mockToken) {
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionToken := fromContext(r.Context())
				if !bytes.Equal(sessionToken.data, token[:]) {
					t.Errorf("Got %#v instead of %#v", sessionToken.data, token[:])
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

			sessionMiddleware := middleware{
				key: key[:],
				newToken: func() tokenizer {
					return &mockToken{
						signFn: func(s *mockToken, k []byte) ([]byte, error) {
							if !bytes.Equal(key[:], k) {
								t.Errorf("keys are not equal %#v != %#v", key[:], k)
							}
							return s.data, nil
						},
						parseFn: func(s *mockToken, key []byte, token []byte) error {
							return fmt.Errorf("invalid token")
						},
						initFn: func(s *mockToken) {
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := sessionMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionToken := fromContext(r.Context())
				if !bytes.Equal(sessionToken.data, token[:]) {
					t.Errorf("Got %#v instead of %#v", sessionToken.data, token[:])
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

			middleware := middleware{
				key: key[:],
				newToken: func() tokenizer {
					return &mockToken{
						signFn: func(s *mockToken, k []byte) ([]byte, error) {
							t.Error("this func does not need to be called")
							return s.data, nil
						},
						parseFn: func(s *mockToken, key []byte, got []byte) error {
							if !bytes.Equal(token[:], got) {
								t.Errorf("Got %#v instead of %#v.", got, token[:])
							}
							s.data = got
							return nil
						},
						initFn: func(s *mockToken) {
							t.Error("this func does not need to be called")
							gen.Read(token[:])
							s.data = token[:]
						},
					}
				},
			}

			handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionToken := fromContext(r.Context())
				if !bytes.Equal(sessionToken.data, token[:]) {
					t.Errorf("Got %#v instead of %#v", sessionToken.data, token[:])
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
