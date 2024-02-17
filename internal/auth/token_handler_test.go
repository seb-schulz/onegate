package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenHandler_GetAndVerifyClient(t *testing.T) {
	newCustomRequest := func(u, p string) *http.Request {
		r := httptest.NewRequest("GET", "/foo", nil)

		r.SetBasicAuth(u, p)
		return r
	}

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
				return &mockClient{"1", "/"}, nil
			}, func(c client) {
				if got := c.ClientID(); got != "1" {
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
				return &mockClient{"1", "/"}, nil
			}, func(c client) {
				if got := c.ClientID(); got != "1" {
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
				return &mockClient{"1", "/"}, nil
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
				return &mockClient{"1", "/"}, nil
			}, func(c client) {
				if got := c.ClientID(); got != "1" {
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
