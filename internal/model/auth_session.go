package model

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type AuthSession struct {
	gorm.Model
	SessionID uint
	Data      webauthn.SessionData `gorm:"serializer:json"`
}
