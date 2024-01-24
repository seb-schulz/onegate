package middleware

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/seb-schulz/onegate/internal/model"
)

func TestParseToken(t *testing.T) {
	key := []byte(".test.")
	ls := loginService{
		key:          key,
		validMethods: []string{"HS256"},
	}

	type testCase struct {
		method             jwt.SigningMethod
		userID             uint
		expiresIn          time.Duration
		expectedParseError error
	}

	for _, tc := range []testCase{
		{jwt.SigningMethodHS256, 1, 10 * time.Second, nil},
		{jwt.SigningMethodHS256, 2, -10 * time.Second, nil},
		{jwt.SigningMethodHS256, 2, -31 * time.Second, jwt.ErrTokenExpired},
		{jwt.SigningMethodHS512, 3, 10 * time.Second, jwt.ErrTokenSignatureInvalid},
	} {
		sToken, _ := jwt.NewWithClaims(tc.method, &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tc.expiresIn)),
			ID:        "abcd",
			Subject:   fmt.Sprintf("%x", tc.userID),
		}).SignedString(key)

		got, err := ls.parseToken(sToken)
		if tc.expectedParseError == nil && err != nil {
			t.Fatalf("cannot parse token: %v", err)
		} else if tc.expectedParseError != nil && errors.Is(tc.expectedParseError, err) {

			t.Fatalf("Got error %v instead of %v", err, tc.expectedParseError)
		}

		if gotID, _ := got.Claims.(*loginClaims).UserID(); gotID != tc.userID {
			t.Fatalf("Got user ID %v instead of %v", gotID, tc.userID)
		}
	}
}

func FuzzGetLoginUrl2(f *testing.F) {
	u, _ := url.Parse("https://example.com/login")

	ls := loginService{
		key:          []byte("abcd"),
		validMethods: []string{"HS256"},
		baseUrl:      *u,
	}

	parseUrl := func(t *testing.T, url *url.URL) error {
		pathChunks := strings.Split(url.Path, "/")
		if pathChunks[1] != "login" {
			t.Fatalf("Url does not start with login: %#v", pathChunks[1])
		}
		_, err := ls.parseToken(pathChunks[len(pathChunks)-1])
		return err
	}

	for i := 0; i < 100; i++ {
		f.Add(uint(rand.Int()))
	}
	f.Fuzz(func(t *testing.T, id uint) {
		out, err := ls.GetLoginUrl(id, time.Second)
		if err != nil {
			t.Fatalf("URL for %v: failed: %v", id, err)
		}

		if err := parseUrl(t, out); err != nil {
			t.Fatalf("cannot parse URL %v: %v\n", *out, err)
		}
	})
}

func TestLoginServiceHandler(t *testing.T) {
	newUrl := func(ls *loginService, userID uint, expiresIn time.Duration) string {
		url, _ := ls.GetLoginUrl(userID, expiresIn)
		return fmt.Sprint(url)
	}
	newRequest := func(url string) *http.Request {
		return httptest.NewRequest("GET", url, nil)
	}
	u, _ := url.Parse("https://example.com/login")

	ls := loginService{
		key:          []byte("abcd"),
		validMethods: []string{"HS256"},
		baseUrl:      *u,
	}

	handler := chi.NewRouter()
	handler.Get("/login/{token}", ls.Handler)

	for _, tc := range []struct {
		expectLogin func(model.LoginOpt)
		req         *http.Request
	}{
		{func(opt model.LoginOpt) {
			t.Error("func should not be called because of an invalid token")
		}, newRequest("http://example.com/login/invalid")},
		{func(lo model.LoginOpt) {
			if lo.UserID == nil || *lo.UserID != 1 {
				t.Errorf("cannot login with valid user ID")
			}
		}, newRequest(newUrl(&ls, 1, 10*time.Second))},
	} {
		ls.loginUser = func(ctx context.Context, opt model.LoginOpt) error {
			tc.expectLogin(opt)
			return nil
		}

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, tc.req)

		resp := w.Result()

		if resp.StatusCode != http.StatusSeeOther {
			t.Errorf("expected status code %v but got %v", http.StatusSeeOther, resp.StatusCode)
		}
	}

}
