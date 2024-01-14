package graph

import (
	"context"
	"fmt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
)

func (r *mutationResolver) beginRegistration(ctx context.Context, user webauthn.User) (*protocol.CredentialCreation, error) {
	options, webauthn_session, err := r.WebAuthn.BeginRegistration(user, webauthn.WithCredentialParameters([]protocol.CredentialParameter{{Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgES256}, {Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgRS256}}), webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementPreferred))
	if err != nil {
		return nil, err
	}

	if _, err := sessionmgr.ContextWithToken[*model.AuthSession](ctx, model.CreateAuthSession(webauthn_session)); err != nil {
		return nil, fmt.Errorf("cannot start registration: %v", err)
	}
	return options, nil
}
