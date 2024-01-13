package model

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthSession struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      datatypes.JSONType[webauthn.SessionData]
}

func (a AuthSession) Value() webauthn.SessionData {
	return a.Data.Data()
}

func CreateAuthSession(data *webauthn.SessionData) func(*gorm.DB, *sessionmgr.Token) (*AuthSession, error) {
	return func(tx *gorm.DB, token *sessionmgr.Token) (*AuthSession, error) {
		authSession := AuthSession{ID: token.UUID, Data: datatypes.NewJSONType(*data)}
		if r := tx.Save(&authSession); r.Error != nil {
			return nil, r.Error
		}
		return &authSession, nil
	}
}

func FirstAuthSession(tx *gorm.DB, token *sessionmgr.Token) (*AuthSession, error) {
	auth_session := AuthSession{ID: token.UUID}
	if result := tx.Order("created_at DESC").First(&auth_session); result.Error != nil {
		return nil, result.Error
	}
	return &auth_session, nil
}
