package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	dbmodel "github.com/seb-schulz/onegate/internal/model"
	"gorm.io/gorm"
)

// ID is the resolver for the ID field.
func (r *credentialResolver) ID(ctx context.Context, obj *dbmodel.Credential) (string, error) {
	return fmt.Sprintf("%d", obj.ID), nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, name string) (*protocol.CredentialCreation, error) {
	session := mustSessionFromContext(ctx)
	time.Sleep(2 * time.Second)

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

	return r.beginRegistration(user, session.ID)
}

// InitCredential is the resolver for the initCredential field.
func (r *mutationResolver) InitCredential(ctx context.Context) (*protocol.CredentialCreation, error) {
	session := mustSessionFromContext(ctx)
	time.Sleep(2 * time.Second)

	if session.UserID == nil {
		return nil, fmt.Errorf("user not logged in")
	}

	return r.beginRegistration(session.User, session.ID)
}

// AddPasskey is the resolver for the addPasskey field.
func (r *mutationResolver) AddCredential(ctx context.Context, body string) (bool, error) {
	session := mustSessionFromContext(ctx)

	if session.UserID == nil {
		return false, fmt.Errorf("no user logged in")
	}

	auth_session := dbmodel.AuthSession{SessionID: session.ID}
	if result := r.DB.Order("updated_at DESC").First(&auth_session, "session_id = ?", session.ID); result.Error != nil {
		return false, fmt.Errorf("no registration session found")
	}
	webauthn_session := auth_session.Data

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(strings.NewReader(body))
	if err != nil {
		return false, err
	}
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

// UpdateCredential is the resolver for the updateCredential field.
func (r *mutationResolver) UpdateCredential(ctx context.Context, id string, description *string) (*dbmodel.Credential, error) {
	if description == nil {
		return nil, fmt.Errorf("no mutation required")
	}

	if len(*description) > 255 {
		return nil, fmt.Errorf("length of description must be less or equal 255 characters")
	}
	session := mustSessionFromContext(ctx)

	if session.UserID == nil {
		return nil, fmt.Errorf("user not logged in")
	}

	cred, err := dbmodel.CredentialByUserID(r.DB, *session.UserID, id)
	if err != nil {
		return nil, err
	}
	cred.Description = *description
	r.DB.Save(&cred)

	return cred, nil
}

// RemoveCredential is the resolver for the removeCredential field.
func (r *mutationResolver) RemoveCredential(ctx context.Context, id string) (bool, error) {
	session := mustSessionFromContext(ctx)

	if session.UserID == nil {
		return false, fmt.Errorf("user not logged in")
	}

	if len(session.User.Credentials) <= 1 {
		return false, fmt.Errorf("cannot delete last remaining credential")
	}

	cred, err := dbmodel.CredentialByUserID(r.DB, *session.UserID, id)
	if err != nil {
		return false, err
	}

	r.DB.Delete(&cred)
	return true, nil
}

// BeginLogin is the resolver for the beginLogin field.
func (r *mutationResolver) BeginLogin(ctx context.Context) (*protocol.CredentialAssertion, error) {
	session := mustSessionFromContext(ctx)

	time.Sleep(2 * time.Second)

	if session.UserID != nil {
		return nil, fmt.Errorf("user is logged-in")
	}

	cred, webauthn_session, err := r.WebAuthn.BeginDiscoverableLogin()
	if err != nil {
		return nil, err
	}
	r.DB.Create(&dbmodel.AuthSession{SessionID: session.ID, Data: *webauthn_session})

	return cred, nil
}

// ValidateLogin is the resolver for the validateLogin field.
func (r *mutationResolver) ValidateLogin(ctx context.Context, body string) (bool, error) {
	session := mustSessionFromContext(ctx)

	time.Sleep(2 * time.Second)

	if session.UserID != nil {
		return false, fmt.Errorf("user is logged-in")
	}

	auth_session := dbmodel.AuthSession{SessionID: session.ID}
	if result := r.DB.Order("updated_at DESC").First(&auth_session, "session_id = ?", session.ID); result.Error != nil {
		return false, fmt.Errorf("login failed")
	}
	webauthn_session := auth_session.Data

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(strings.NewReader(body))
	if err != nil {
		return false, fmt.Errorf("login failed")
	}

	user := dbmodel.User{}
	db_cred := dbmodel.Credential{}

	cred, err := r.WebAuthn.ValidateDiscoverableLogin(func(rawID, userHandle []byte) (webauthn.User, error) {
		if result := r.DB.Preload("Credentials").First(&user, "authn_id = ?", userHandle); result.Error != nil {
			return nil, fmt.Errorf("login failed")
		}

		for _, c := range user.Credentials {
			if bytes.Equal(c.Data.ID, rawID) {
				db_cred = c
				return &user, nil
			}
		}

		return nil, fmt.Errorf("login failed")
	}, webauthn_session, parsedResponse)
	if err != nil {
		return false, fmt.Errorf("login failed")
	}

	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		db_cred.Data = *cred
		session.User = &user
		tx.Save(&db_cred)
		tx.Save(&session)
		tx.Delete(&auth_session)
		return nil
	}); err != nil {
		panic(err)
	}

	return true, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*dbmodel.User, error) {
	session := mustSessionFromContext(ctx)

	if session.UserID == nil {
		return nil, nil // ignore error to dedect logged-out scenario
	}
	r.DB.Model(&session.User).First(&session.User)

	return session.User, nil
}

// Credentials is the resolver for the credentials field.
func (r *queryResolver) Credentials(ctx context.Context) ([]*dbmodel.Credential, error) {
	session := mustSessionFromContext(ctx)

	if session.UserID == nil {
		return nil, nil // ignore error to dedect logged-out scenario
	}

	creds := []*dbmodel.Credential{}
	if result := r.DB.Where("user_id = ?", session.UserID).Find(&creds); result.Error != nil {
		return nil, result.Error
	}

	return creds, nil
}

// Credential returns CredentialResolver implementation.
func (r *Resolver) Credential() CredentialResolver { return &credentialResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type credentialResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
