package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/model"
	"gorm.io/gorm"

	"golang.org/x/oauth2"
)

type mockClient struct {
	c uuid.UUID
	r string
}

func (mc *mockClient) ClientID() uuid.UUID {
	return mc.c
}

func (mc *mockClient) VerifyClientSecret(s string) error {
	if s != "secret" {
		return fmt.Errorf("secret does not match")
	}
	return nil
}

func (mc *mockClient) RedirectURI() string {
	return mc.r
}

type mockAuthorization struct {
	Authorization
	client mockClient
	userID *uint
}

func (ma *mockAuthorization) Code() string {
	return "mno"
}

func (ma *mockAuthorization) RedirectURI() string {
	return ma.client.RedirectURI()
}

func (ma *mockAuthorization) SetUserID(ctx context.Context, userID uint) error {
	ma.userID = &userID
	return nil
}

func TestAuthCodeFlow(t *testing.T) {
	var (
		code         string
		expectedBody = "Done \\o/"
	)

	client_ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code = r.FormValue("code")
		if code == "" {
			t.Errorf("expected non-empty authorization code")
		}

		state := r.FormValue("state")
		if state != "state" {
			t.Errorf("expected 'state' as state value but got: %#v", state)
		}
		fmt.Fprint(w, expectedBody)
	}))
	defer client_ts.Close()

	mockClient := mockClient{
		uuid.MustParse("2e532bfa50a44f1c84aa5af13fa4612d"),
		client_ts.URL,
	}

	mockUser := model.User{
		Model: gorm.Model{ID: 1},
	}

	mock := struct {
		currentAuthorization *mockAuthorization
	}{}

	clientByClientID := func(ctx context.Context, clientID string) (client, error) {
		return &mockClient, nil
	}

	route := chi.NewRouter()
	authorizationRequestHandler := &authorizationRequestHandler{
		clientByClientID: clientByClientID,
		loginUrl:         url.URL{Path: "/login"},
		createAuthorization: func(ctx context.Context, client client, state, codeChallenge string) error {
			mock.currentAuthorization = &mockAuthorization{
				Authorization{
					InternalState:         state,
					InternalCodeChallenge: codeChallenge,
					InternalClientID:      client.ClientID(),
				},
				mockClient,
				nil,
			}
			return nil
		},
	}
	route.Get("/auth", authorizationRequestHandler.ServeHTTP)
	route.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Instead mock, it should call assignUserToAuthorization(...)
		// sessionToken := sessionmgr.Token{UUID: uuid.New()}
		uID := uint(1)
		mock.currentAuthorization.InternalUserID = &uID
		// mock.currentAuthorization.User = &model.User{Model: gorm.Model{ID: 1}}
		http.Redirect(w, r, "/callback", http.StatusFound)

	})

	callbackRedirectHandler := &callbackRedirectHandler{
		currentAuthorization: func(ctx context.Context) (authorization, error) {
			return mock.currentAuthorization, nil
		},
		currentUser: func(ctx context.Context) *model.User {
			return &mockUser
		},
	}
	route.Get("/callback", callbackRedirectHandler.ServeHTTP)

	tokenHandler := &tokenHandler{
		clientByClientID: clientByClientID,
		authorizationByCode: func(ctx context.Context, code string) (authorization, error) {
			return mock.currentAuthorization, nil
		},
		deleteAuthorization: func(ctx context.Context, a authorization) error {
			return nil
		},
	}
	route.Post("/token", tokenHandler.ServeHTTP)

	ts := httptest.NewServer(route)
	defer ts.Close()

	conf := &oauth2.Config{
		ClientID:     "123",
		ClientSecret: "secret",
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%v/auth", ts.URL),
			TokenURL:  fmt.Sprintf("%v/token", ts.URL),
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()
	authUrl := conf.AuthCodeURL("state", oauth2.S256ChallengeOption(verifier))

	res, err := http.Get(authUrl)
	if err != nil {
		t.Errorf("cannot get authorization token: %v", err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if string(greeting) != expectedBody {
		t.Errorf("Got %#v instead of %#v", string(greeting), expectedBody)
	}

	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tok)

	if mock.currentAuthorization.userID == nil {
		t.Errorf("user ID was not set")
	} else if *mock.currentAuthorization.userID != mockUser.ID {
		t.Errorf("expected user ID %v but got %v", mockUser.ID, *mock.currentAuthorization.userID)
	}
	// t.Log(tok.Extra("id_token"))

	// client := conf.Client(ctx, tok)
	// client.Get("...")
	// t.FailNow()
}
