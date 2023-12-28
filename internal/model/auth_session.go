package model

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthSession struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	SessionID uint
	Data      datatypes.JSONType[webauthn.SessionData]
}

func (a AuthSession) Value() webauthn.SessionData {
	return a.Data.Data()
}

func CreateAuthSession(tx *gorm.DB, sID uint, data webauthn.SessionData) error {
	if r := tx.Create(&AuthSession{SessionID: sID, Data: datatypes.NewJSONType(data)}); r.Error != nil {
		return r.Error
	}

	return nil
}

func FirstAuthSessionBySession(tx *gorm.DB, sID uint) (*AuthSession, error) {
	auth_session := AuthSession{}
	if result := tx.Order("created_at DESC").First(&auth_session, "session_id = ?", sID); result.Error != nil {
		return nil, result.Error
	}
	return &auth_session, nil
}
