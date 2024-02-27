package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
)

type authorization interface {
	ClientID() uuid.UUID
	UserID() uint
	State() string
	Code() string
	CodeChallenge() string
	redirecter
	SetUserID(context.Context, uint) error
}

type Authorization struct {
	ID                    uint `gorm:"primarykey"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	InternalClientID      uuid.UUID   `gorm:"column:client_id;type:VARCHAR;size:191;not null"`
	Client                Client      `gorm:"foreignKey:InternalClientID"`
	InternalUserID        *uint       `gorm:"column:user_id"`
	User                  *model.User `gorm:"foreignKey:InternalUserID"`
	InternalState         string      `gorm:"column:state"`
	InternalCode          []byte      `gorm:"column:code;type:BLOB(16)"`
	InternalCodeChallenge string      `gorm:"column:code_challenge;type:BLOB(16)"`
	SessionID             uuid.UUID   `gorm:"column:session_id;type:VARCHAR(191);not null"`
}

func (a *Authorization) ClientID() uuid.UUID {
	return a.InternalClientID
}

func (a *Authorization) UserID() uint {
	if a.InternalUserID == nil {
		return 0
	}
	return *a.InternalUserID
}

func (a *Authorization) State() string {
	return a.InternalState
}

func (a *Authorization) Code() string {
	return base64.URLEncoding.EncodeToString(a.InternalCode)
}

func (a *Authorization) CodeChallenge() string {
	return a.InternalCodeChallenge
}

func (a *Authorization) RedirectURI() string {
	return a.Client.RedirectURI()
}

func (a *Authorization) IDStr() string {
	return fmt.Sprint(a.ID)
}

func (a *Authorization) SetUserID(ctx context.Context, userID uint) error {
	r := database.FromContext(ctx).Model(a).Update("user_id", userID)
	if r.Error != nil {
		return fmt.Errorf("cannot update authorization: %w", r.Error)
	}

	return nil
}

func createAuthorization(ctx context.Context, client client, state, codeChallenge string) error {

	if state == "" {
		return fmt.Errorf("state must not be empty")
	}

	if codeChallenge == "" {
		return fmt.Errorf("code challenge must not be empty")
	}

	code := make([]byte, 16)
	if err := readRand(code); err != nil {
		panic("cannot generate code")

	}

	authReq := Authorization{
		InternalClientID:      client.ClientID(),
		InternalState:         state,
		InternalCodeChallenge: codeChallenge,
		InternalCode:          code,
		SessionID:             sessionmgr.FromContext(ctx).UUID,
	}

	r := database.FromContext(ctx).Create(&authReq)
	if r.Error != nil {
		return fmt.Errorf("cannot create authorization: %v", r.Error)
	}

	return nil
}

func FirstAuthorization(ctx context.Context) (*Authorization, error) {
	sessionID := sessionmgr.FromContext(ctx).UUID
	authReq := Authorization{}
	r := database.FromContext(ctx).Preload("Client").Where("session_id = ?", sessionID).First(&authReq)

	if r.Error != nil {
		return nil, fmt.Errorf("cannot update authorization: %v", r.Error)
	}

	return &authReq, nil
}

func authorizationByCode(ctx context.Context, code string) (authorization, error) {
	decCode, err := base64.URLEncoding.DecodeString(code)
	if err != nil {
		return nil, err
	}

	authReq := Authorization{}
	r := database.FromContext(ctx).Where("code = ?", decCode).First(&authReq)

	if r.Error != nil {
		return nil, fmt.Errorf("cannot get authorization: %w", r.Error)
	}

	return &authReq, nil
}
