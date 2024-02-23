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
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
	"gorm.io/gorm"

	"golang.org/x/oauth2"
)

type mockClient struct{ c, r string }

func (mc *mockClient) ClientID() string {
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

func (ma *mockAuthorization) Exists() bool {
	return true
}

func (ma *mockAuthorization) ClientID() string {
	return "123"
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

func (ma *mockAuthorization) IDStr() string {
	return ma.ClientID()
}

type mockAuthorizationMgr struct {
	*sessionmgr.StorageManager[*mockAuthorization]
	currentAuthorization *mockAuthorization
}

func (auth *mockAuthorizationMgr) create(ctx context.Context, client client, state, codeChallenge string) error {
	auth.currentAuthorization = &mockAuthorization{state: state, codeChallenge: codeChallenge, redirectURI: client.RedirectURI()}
	return nil
}

func (auth *mockAuthorizationMgr) updateUserID(ctx context.Context, userID uint) error {
	auth.currentAuthorization.userID = 1
	return nil
}

func (auth *mockAuthorizationMgr) byCode(ctx context.Context, code string) (authorization, error) {
	return auth.currentAuthorization, nil
}

func (auth *mockAuthorizationMgr) fromContext(ctx context.Context) authorization {
	return auth.currentAuthorization
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

	mock := &mockAuthorizationMgr{}
	mock.StorageManager = sessionmgr.NewStorage[*mockAuthorization]("authorization", func(ctx context.Context) (*mockAuthorization, error) {
		return mock.currentAuthorization, nil
	})

	clientByClientID := func(ctx context.Context, clientID string) (client, error) {
		return &mockClient{"123", client_ts.URL}, nil
	}

	route := chi.NewRouter()
	authorizationRequestHandler := &authorizationRequestHandler{
		clientByClientID: clientByClientID,
		loginUrl:         url.URL{Path: "/callback"},
		authorizationMgr: mock,
	}
	route.Get("/auth", authorizationRequestHandler.ServeHTTP)

	authorizationResponseHandler := &authorizationResponseHandler{
		authorizationMgr: mock,
		currentUserFromContext: func(ctx context.Context) *model.User {
			return &model.User{Model: gorm.Model{ID: 1}}
		},
	}
	route.With(mock.Handler).Get("/callback", authorizationResponseHandler.ServeHTTP)

	tokenHandler := &tokenHandler{clientByClientID: clientByClientID, authorizationMgr: mock}
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
