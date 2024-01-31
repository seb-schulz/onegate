package server

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
	"gorm.io/gorm"
)

func TestRedirectWhenLoggedOut(t *testing.T) {
	type testCase struct {
		expectedStatus int
		fromCtx        func(context.Context) *model.User
	}

	_orig := defaultUserFromContext
	defer func() {
		defaultUserFromContext = _orig
	}()

	for _, tc := range []testCase{
		{http.StatusOK, func(context.Context) *model.User {
			return nil
		}},
		{http.StatusSeeOther, func(context.Context) *model.User {
			return &model.User{Model: gorm.Model{ID: 1}}
		}},
	} {
		defaultUserFromContext = tc.fromCtx
		handler := redirectWhenLoggedOut(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Ok")
		}))

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("GET", "/foo", nil))

		resp := w.Result()
		if resp.StatusCode != tc.expectedStatus {
			t.Errorf("Got status code %#v and expected %#v", resp.StatusCode, tc.expectedStatus)
		}

	}
}

func TestTokenBasedLoginServiceParseToken(t *testing.T) {
	key := []byte(".test.")
	ls := tokenBasedLoginService{
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

func FuzzTokenBasedLoginServiceGetLoginUrl(f *testing.F) {
	u, _ := url.Parse("https://example.com/login")

	ls := tokenBasedLoginService{
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
		out, err := ls.getLoginUrl(id, time.Second)
		if err != nil {
			t.Fatalf("URL for %v: failed: %v", id, err)
		}

		if err := parseUrl(t, out); err != nil {
			t.Fatalf("cannot parse URL %v: %v\n", *out, err)
		}
	})
}

func TestTokenBasedLoginServiceHandler(t *testing.T) {
	newUrl := func(ls *tokenBasedLoginService, userID uint, expiresIn time.Duration) string {
		url, _ := ls.getLoginUrl(userID, expiresIn)
		return fmt.Sprint(url)
	}
	newRequest := func(url string) *http.Request {
		return httptest.NewRequest("GET", url, nil)
	}
	ls := tokenBasedLoginService{
		key:          []byte("abcd"),
		validMethods: []string{"HS256"},
		baseUrl:      url.URL{Scheme: "https", Host: "example.com", Path: "/login"},
	}

	handler := chi.NewRouter()
	handler.Get("/login/{token}", ls.handler)

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
