package auth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
	"gorm.io/gorm"

	"golang.org/x/oauth2"
)

type tokenHandler struct {
	clientByClientID clientByClientIDFn
	authorizationMgr interface {
		byCode(ctx context.Context, code string) (authorization, error)
	}
	ClientSecretVerifier
}

func (th *tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client, err := th.getAndVerifyClient(r)
	if err != nil {
		log.Printf("cannot verify client: %v", err)
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}

	if err := th.checkGrantType(r); err != nil {
		http.Error(w, errors.Descriptions[err], errors.StatusCodes[err])
		return
	}

	authReq, err := th.authorizationMgr.byCode(r.Context(), r.FormValue("code"))
	if err != nil {
		log.Printf("authorization not found: %v", err)
		http.Error(w, "not implemented yet", http.StatusNotImplemented)
		return
	}

	if err := th.checkCodeChallenge(r, authReq); err != nil {
		log.Println("missmach between authorization and client")
		http.Error(w, "not implemented yet", http.StatusNotImplemented)
		return
	}

	if authReq.ClientID() != client.ClientID() {
		log.Println("missmach between authorization and client")
		http.Error(w, "not implemented yet", http.StatusNotImplemented)
		return
	}

	b, err := json.Marshal(struct {
		AccessToken string `json:"access_token,omitempty"`
		TokenType   string `json:"token_type,omitempty"`
		ExpiresIn   int    `json:"expires_in,omitempty"`
		IDToken     string `json:"id_token,omitempty"`
	}{"xyz123", "Bearer", 10 * 60, "abc"})
	if err != nil {
		log.Printf("cannot generate token: %v", err)
		http.Error(w, "not implemented yet", http.StatusNotImplemented)
		return
	}

	log.Printf("Form params: %#v", r.Form)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (th *tokenHandler) checkCodeChallenge(r *http.Request, auth authorization) error {
	cv := r.FormValue("code_verifier")
	if subtle.ConstantTimeCompare([]byte(auth.CodeChallenge()), []byte(oauth2.S256ChallengeFromVerifier(cv))) != 1 {
		return fmt.Errorf("missmatch of code challenge")
	}
	return nil
}

func (th *tokenHandler) checkGrantType(r *http.Request) error {
	if r.FormValue("grant_type") != "authorization_code" {
		return errors.ErrInvalidGrant
	}
	return nil
}

func (th *tokenHandler) getAndVerifyClient(r *http.Request) (client, error) {
	var clientID, secret string

	clientID, secret, ok := r.BasicAuth()
	if !ok {
		clientID = r.FormValue("client_id")
		secret = r.FormValue("client_secret")
	}

	client, err := th.clientByClientID(r.Context(), clientID)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch client: %v", err)
	}

	if err := client.VerifyClientSecret(secret); err != nil {
		return nil, err
	}
	return client, nil
}

type mockClient struct{ r string }

func (mc *mockClient) ClientID() string {
	return "123"
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

func (auth *mockAuthorizationMgr) FromContext(ctx context.Context) authorization {
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
		return &mockClient{client_ts.URL}, nil
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
