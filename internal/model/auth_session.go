package model

import (
	"context"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"
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

func CreateAuthSession(ctx context.Context, data *webauthn.SessionData) error {
	_, err := database.Transaction(ctx, func(tx *gorm.DB) (bool, error) {
		authSession := AuthSession{
			ID:   sessionmgr.FromContext(ctx).UUID,
			Data: datatypes.NewJSONType(*data),
		}
		if r := tx.Save(&authSession); r.Error != nil {
			return false, r.Error
		}
		return true, nil
	})
	return err
}

func FirstAuthSession(ctx context.Context) (*AuthSession, error) {
	auth_session := AuthSession{ID: sessionmgr.FromContext(ctx).UUID}
	if result := database.FromContext(ctx).Order("created_at DESC").First(&auth_session); result.Error != nil {
		return nil, result.Error
	}
	return &auth_session, nil
}
