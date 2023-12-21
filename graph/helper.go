package graph

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/model"
)

func mustRandomEncodedBytes(len int) string {
	r := make([]byte, len)

	_, err := rand.Read(r)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(r)
}

func mustSessionFromContext(ctx context.Context) *model.Session {
	session := middleware.SessionFromContext(ctx)
	if session == nil {
		panic("session is missing")
	}
	return session
}

func (r *mutationResolver) beginRegistration(user webauthn.User, sessionID uint) (*protocol.CredentialCreation, error) {
	options, webauthn_session, err := r.WebAuthn.BeginRegistration(user, webauthn.WithCredentialParameters([]protocol.CredentialParameter{{Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgES256}, {Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgRS256}}), webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementPreferred))
	if err != nil {
		return nil, err
	}
	r.DB.Create(&model.AuthSession{SessionID: sessionID, Data: *webauthn_session})

	return options, nil
}
