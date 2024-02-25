package auth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/errors"
	"golang.org/x/oauth2"
)

type AccessTokenResponds struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	IDToken     string `json:"id_token,omitempty"`
}

type tokenHandler struct {
	clientByClientID    clientByClientIDFn
	authorizationByCode func(ctx context.Context, code string) (authorization, error)
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

	authReq, err := th.authorizationByCode(r.Context(), r.FormValue("code"))
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

	b, err := json.Marshal(AccessTokenResponds{"xyz123", "Bearer", 10 * 60, "abc"})
	if err != nil {
		log.Printf("cannot generate token: %v", err)
		http.Error(w, "not implemented yet", http.StatusNotImplemented)
		return
	}

	// log.Printf("Form params: %#v", r.Form)

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