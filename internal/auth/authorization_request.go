package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-oauth2/oauth2/v4/errors"
)

type authorizationRequestHandler struct {
	clientByClientID clientByClientIDFn
	authorizationMgr interface {
		create(ctx context.Context, client client, state string, codeChallenge string) error
	}
	loginUrl url.URL
}

func (auth authorizationRequestHandler) checkResponseType(response_type string) error {
	if response_type != "code" {
		return errors.ErrUnsupportedResponseType
	}
	return nil
}

func (auth authorizationRequestHandler) checkMethod(r *http.Request) error {
	if !(r.Method == "GET" || r.Method == "POST") {
		return errors.ErrInvalidRequest
	}

	return nil
}

func (auth authorizationRequestHandler) checkCodeChallengeMethod(r *http.Request) error {
	if r.FormValue("code_challenge_method") != "S256" {
		return errors.ErrInvalidRequest
	}

	return nil
}

func (auth authorizationRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Query params: %#v", r.URL.Query())
	if err := auth.checkMethod(r); err != nil {
		http.Error(w, errors.Descriptions[err], errors.StatusCodes[err])
		return
	}

	if err := auth.checkResponseType(r.FormValue("response_type")); err != nil {
		http.Error(w, errors.Descriptions[err], errors.StatusCodes[err])
		return
	}

	if err := auth.checkCodeChallengeMethod(r); err != nil {
		http.Error(w, errors.Descriptions[err], errors.StatusCodes[err])
		return
	}

	client, err := auth.clientByClientID(r.Context(), r.FormValue("client_id"))
	if err != nil {
		http.Error(w, errors.Descriptions[errors.ErrInvalidRequest], errors.StatusCodes[errors.ErrInvalidRequest])
		return
	}

	if auth.authorizationMgr.create(r.Context(), client, r.FormValue("state"), r.FormValue("code_challenge")); err != nil {
		http.Error(w, errors.Descriptions[errors.ErrInvalidRequest], errors.StatusCodes[errors.ErrInvalidRequest])
		return
	}

	http.Redirect(w, r, fmt.Sprint(&auth.loginUrl), http.StatusSeeOther)
}
