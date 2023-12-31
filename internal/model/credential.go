package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	UserID      uint
	User        User
	Description string
	LastLogin   *time.Time
	Data        webauthn.Credential `gorm:"serializer:json"`
}

func CredentialByUserID(db *gorm.DB, userID uint, id string) (*Credential, error) {
	cred := Credential{}
	r := db.Where("user_id = ? AND id = ?", userID, id).First(&cred)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("credential not found")
	}
	return &cred, nil
}

func CountCredentialByUserID(db *gorm.DB, userID uint) int {
	var c int64
	db.Model(&Credential{}).Where("user_id = ?", userID).Count(&c)
	return int(c)

}
