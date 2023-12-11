package model

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	UserID int
	User   User
	Data   webauthn.Credential `gorm:"serializer:json"`
}
