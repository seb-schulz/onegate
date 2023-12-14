package model

import (
	"fmt"
	"log"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AuthnID     []byte `gorm:"type:BLOB(16);default:RANDOM_BYTES(16);not null"`
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

func GetUserByRawID(rawID, userHandle []byte) (webauthn.User, error) {
	log.Println("GetUserByRawID", rawID, userHandle)
	return nil, fmt.Errorf("user not found")
}
