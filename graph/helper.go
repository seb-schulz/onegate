package graph

import (
	"context"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/model"
)

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
	model.CreateAuthSession(r.DB, sessionID, *webauthn_session)
	return options, nil
}
