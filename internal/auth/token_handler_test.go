package auth

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func TestTokenHandler_GetAndVerifyClient(t *testing.T) {
	newCustomRequest := func(u, p string) *http.Request {
		r := httptest.NewRequest("GET", "/foo", nil)

		r.SetBasicAuth(u, p)
		return r
	}

	testClientID := uuid.MustParse("2e532bfa50a44f1c84aa5af13fa4612d")

	for _, tc := range []struct {
		inputRequest *http.Request
		getClient    clientByClientIDFn
		checkClient  func(client)
		checkError   func(error)
	}{
		{
			httptest.NewRequest("GET", "/foo", nil), func(ctx context.Context, clientID string) (client, error) {
				return nil, fmt.Errorf("failed")
			}, func(c client) {
				t.Error("should not be called due to error")
			}, func(err error) {
				if err == nil {
					t.Error("expected error but got nil")
				}
			},
		}, {
			newCustomRequest("1", "invalid"), func(ctx context.Context, clientID string) (client, error) {
				if clientID != "1" {
					t.Errorf("expected client ID 1 but got: %#v", clientID)
				}
				return &mockClient{testClientID, "/"}, nil
			}, func(c client) {
				if got := c.ClientID(); got != testClientID {
					t.Errorf("expected client ID 1 but got: %#v", got)
				}
			}, func(err error) {
				if err == nil {
					t.Error("expected error but got nil")
				}
			},
		}, {
			newCustomRequest("1", "secret"), func(ctx context.Context, clientID string) (client, error) {
				if clientID != "1" {
					t.Errorf("expected client ID 1 but got: %#v", clientID)
				}
				return &mockClient{testClientID, "/"}, nil
			}, func(c client) {
				if got := c.ClientID(); got != testClientID {
					t.Errorf("expected client ID 1 but got: %#v", got)
				}
			}, func(err error) {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			},
		}, {
			httptest.NewRequest("GET", "/foo?client_id=1&client_secret=invalid", nil), func(ctx context.Context, clientID string) (client, error) {
				if clientID != "1" {
					t.Errorf("expected client ID 1 but got: %#v", clientID)
				}
				return &mockClient{testClientID, "/"}, nil
			}, func(c client) {
				t.Error("should not be called due to error")
			}, func(err error) {
				if err == nil {
					t.Error("expected error but got nil")
				}
			},
		}, {
			httptest.NewRequest("GET", "/foo?client_id=1&client_secret=secret", nil), func(ctx context.Context, clientID string) (client, error) {
				if clientID != "1" {
					t.Errorf("expected client ID 1 but got: %#v", clientID)
				}
				return &mockClient{testClientID, "/"}, nil
			}, func(c client) {
				if got := c.ClientID(); got != testClientID {
					t.Errorf("expected client ID 1 but got: %#v", got)
				}
			}, func(err error) {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			},
		},
	} {
		handler := &tokenHandler{
			clientByClientID: tc.getClient,
		}
		c, err := handler.getAndVerifyClient(tc.inputRequest)
		if err == nil {
			tc.checkClient(c)
		}
		tc.checkError(err)
	}
}

func TestTokenHandler_CheckGrantType(t *testing.T) {
	for _, tc := range []struct {
		inputRequest  *http.Request
		expectedError error
	}{
		{
			httptest.NewRequest("GET", "/foo?grant_type=authorization_code", nil), nil,
		},
		{
			httptest.NewRequest("GET", "/foo?grant_type=invalid_type", nil), errors.ErrInvalidGrant,
		},
		{
			httptest.NewRequest("GET", "/foo", nil), errors.ErrInvalidGrant,
		},
	} {
		handler := &tokenHandler{}
		err := handler.checkGrantType(tc.inputRequest)
		if err != tc.expectedError {
			t.Errorf("Expected error %v but got %v", tc.expectedError, err)
		}
	}
}

type mockAuthorizationCodeChallenger struct {
	cc string
}

func (mock *mockAuthorizationCodeChallenger) CodeChallenge() string {
	return mock.cc
}

func TestTokenHandler_checkCodeChallenge(t *testing.T) {
	for _, tc := range []struct {
		setup         func() (inputUrl, inputVerifier string)
		expectedError error
	}{
		{
			func() (inputUrl string, inputCodeChallenger string) {
				verifier := oauth2.GenerateVerifier()
				inputCodeChallenger = oauth2.S256ChallengeFromVerifier(verifier)
				inputUrl = fmt.Sprintf("/foo?code_verifier=%s", verifier)
				return
			}, nil,
		},
		{
			func() (inputUrl string, inputCodeChallenger string) {
				origVerifier := oauth2.GenerateVerifier()
				newVerifier := oauth2.GenerateVerifier()
				inputCodeChallenger = oauth2.S256ChallengeFromVerifier(origVerifier)
				inputUrl = fmt.Sprintf("/foo?code_verifier=%s", newVerifier)
				return
			}, errors.ErrInvalidCodeChallenge,
		},
		{
			func() (inputUrl string, inputCodeChallenger string) {
				inputCodeChallenger = oauth2.S256ChallengeFromVerifier(oauth2.GenerateVerifier())
				inputUrl = "/foo"
				return
			}, errors.ErrMissingCodeVerifier,
		},
	} {
		url, cc := tc.setup()
		handler := &tokenHandler{}
		err := handler.checkCodeChallenge(httptest.NewRequest("GET", url, nil), &mockAuthorizationCodeChallenger{cc})

		if err != tc.expectedError {
			t.Errorf("Expected error %#v but got %#v", tc.expectedError, err)
		}
	}
}
func FuzzTokenHandler_checkCodeChallenge(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(rand.Int())
	}

	f.Fuzz(func(t *testing.T, seed int) {
		gen := rand.New(rand.NewSource(int64(seed)))

		newVerifier := make([]byte, 32)
		if _, err := gen.Read(newVerifier); err != nil {
			t.Errorf("failed to setup test: %v", err)
		}

		verifier := oauth2.GenerateVerifier()
		handler := &tokenHandler{}
		err := handler.checkCodeChallenge(httptest.NewRequest("GET", fmt.Sprintf("/foo?code_verifier=%s", url.QueryEscape(string(newVerifier))), nil), &mockAuthorizationCodeChallenger{oauth2.S256ChallengeFromVerifier(verifier)})

		if err == nil {
			t.Errorf("expected error but got no error")
		}
	})
}
