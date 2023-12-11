package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/internal/middleware"
	dbmodel "github.com/seb-schulz/onegate/internal/model"
	"gorm.io/gorm"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, name string) (*protocol.CredentialCreation, error) {
	session := middleware.SessionFromContext(ctx)
	if session == nil {
		return nil, fmt.Errorf("session is missing")
	}

	if session.UserID != nil {
		return nil, fmt.Errorf("currently logged in with an user")
	}

	user := dbmodel.User{Name: name}
	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&user)
		tx.Model(&session).Update("user_id", user.ID)
		return nil
	}); err != nil {
		panic(err)
	}

	options, webauthn_session, err := r.WebAuthn.BeginRegistration(user, webauthn.WithCredentialParameters([]protocol.CredentialParameter{{Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgES256}, {Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgRS256}}))
	if err != nil {
		return nil, err
	}
	r.DB.Create(&dbmodel.AuthSession{SessionID: session.ID, Data: *webauthn_session})

	return options, nil
}

// AddPasskey is the resolver for the addPasskey field.
func (r *mutationResolver) AddPasskey(ctx context.Context, body string) (bool, error) {
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(strings.NewReader(body))
	if err != nil {
		return false, err
	}

	session := middleware.SessionFromContext(ctx)
	if session == nil {
		return false, fmt.Errorf("session is missing")
	}

	if session.UserID == nil {
		return false, fmt.Errorf("no user logged in")
	}

	auth_session := dbmodel.AuthSession{SessionID: session.ID}
	if result := r.DB.First(&auth_session, "session_id = ?", session.ID); result.Error != nil {
		return false, fmt.Errorf("no registration session found")
	}
	webauthn_session := auth_session.Data

	cred, err := r.WebAuthn.CreateCredential(session.User, webauthn_session, parsedResponse)
	if err != nil {
		return false, err
	}

	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&dbmodel.Credential{UserID: *session.UserID, Data: *cred})
		tx.Delete(&auth_session)
		return nil
	}); err != nil {
		panic(err)
	}

	return true, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*dbmodel.User, error) {
	session := middleware.SessionFromContext(ctx)
	if session == nil {
		return nil, fmt.Errorf("session is missing")
	}

	if session.UserID == nil {
		return nil, nil
	}

	return session.User, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
