package model

import (
	"crypto/rand"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
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

func (u *User) IDStr() string {
	return fmt.Sprint(u.ID)
}

func CreateUser(name string) func(tx *gorm.DB, token *sessionmgr.Token) (*User, error) {
	return func(tx *gorm.DB, token *sessionmgr.Token) (*User, error) {
		user := User{Name: name}

		if r := tx.Create(&user); r.Error != nil {
			return nil, r.Error
		}

		if r := tx.FirstOrCreate(&Session{
			ID:   token.UUID,
			User: user,
		}); r.Error != nil {
			return nil, r.Error
		}

		return &user, nil
	}
}

func LoginUser(user *User, cred *Credential) func(tx *gorm.DB, token *sessionmgr.Token) (*User, error) {
	return func(tx *gorm.DB, token *sessionmgr.Token) (*User, error) {
		if r := tx.FirstOrCreate(&Session{
			ID:   token.UUID,
			User: *user,
		}); r.Error != nil {
			return nil, r.Error
		}
		if cred != nil {
			tx.Save(cred)
		}

		return user, nil
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
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
