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
	userID        uint
	state         string
	codeChallenge string
	redirectURI   string
}

func (ma *mockAuthorization) ClientID() uuid.UUID {
	return uuid.MustParse("2e532bfa50a44f1c84aa5af13fa4612d")
}

func (ma *mockAuthorization) UserID() uint {
	return ma.userID
}

func (ma *mockAuthorization) State() string {
	return ma.state
}

func (ma *mockAuthorization) Code() string {
	return "mno"
}

func (ma *mockAuthorization) CodeChallenge() string {
	return ma.codeChallenge
}

func (ma *mockAuthorization) RedirectURI() string {
	return ma.redirectURI
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

	mock := struct {
		currentAuthorization *mockAuthorization
	}{}

	clientByClientID := func(ctx context.Context, clientID string) (client, error) {
		return &mockClient{
			uuid.MustParse("2e532bfa50a44f1c84aa5af13fa4612d"),
			client_ts.URL,
		}, nil
	}

	route := chi.NewRouter()
	authorizationRequestHandler := &authorizationRequestHandler{
		clientByClientID: clientByClientID,
		loginUrl:         url.URL{Path: "/login"},
		createAuthorization: func(ctx context.Context, client client, state, codeChallenge string) error {
			mock.currentAuthorization = &mockAuthorization{state: state, codeChallenge: codeChallenge, redirectURI: client.RedirectURI()}
			return nil
		},
	}
	route.Get("/auth", authorizationRequestHandler.ServeHTTP)
	route.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Instead mock, it should call assignUserToAuthorization(...)
		// sessionToken := sessionmgr.Token{UUID: uuid.New()}
		mock.currentAuthorization.userID = 1
		http.Redirect(w, r, "/callback", http.StatusFound)

	})

	// TODO: Reimplement authorizationResponseHandler to something usefull
	route.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		authReq := mock.currentAuthorization
		q := url.Values{}
		q.Add("code", authReq.Code())
		q.Add("state", authReq.State())
		http.Redirect(w, r, fmt.Sprintf("%v?%v", authReq.RedirectURI(), q.Encode()), http.StatusFound)
	})

	tokenHandler := &tokenHandler{
		clientByClientID: clientByClientID,
		authorizationByCode: func(ctx context.Context, code string) (authorization, error) {
			return mock.currentAuthorization, nil
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
	// log.Println(tok.Extra("id_token"))

	// client := conf.Client(ctx, tok)
	// client.Get("...")
	// t.FailNow()
}
