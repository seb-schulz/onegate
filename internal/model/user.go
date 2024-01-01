package model

import (
	"crypto/rand"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// `RANDOM_BYTES` was added with MariaDB 10.10.0
	// AuthnID     []byte `gorm:"type:BLOB(16);default:RANDOM_BYTES(16);not null"`
	AuthnID     []byte `gorm:"type:BLOB(16)"`
	Name        string `gorm:"type:VARCHAR(255);not null"`
	DisplayName string `gorm:"type:VARCHAR(255);not null"`
	Credentials []Credential
}

func (u User) WebAuthnID() []byte {
	return []byte(u.AuthnID)
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	r := make([]webauthn.Credential, len(u.Credentials))
	for i, v := range u.Credentials {
		r[i] = v.Data
	}
	return r
}

func (u User) WebAuthnIcon() string {
	return ""
}

func CreateUser(user *User, session *Session) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		if user == nil || session == nil {
			return fmt.Errorf("user and session must be provided")
		}

		if session.UserID != nil {
			return fmt.Errorf("session assigned to user")
		}

		if r := tx.Create(&user); r.Error != nil {
			return r.Error
		}

		session.User = user
		if r := tx.Save(&session); r.Error != nil {
			return r.Error
		}
		return nil
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.AuthnID) > 0 {
		return nil
	}

	r := make([]byte, 16)
	if _, err := rand.Read(r); err != nil {
		return err
	}

	u.AuthnID = r
	return nil
}
