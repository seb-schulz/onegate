package model

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AuthnID     []byte `gorm:"type:BLOB(16);default:RANDOM_BYTES(16);not null"`
	Name        string `gorm:"type:VARCHAR(255);not null"`
	DisplayName string `gorm:"type:VARCHAR(255);not null"`
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
	return []webauthn.Credential{}
}

func (u User) WebAuthnIcon() string {
	return ""
}

func MarshalCredentialCreation(c protocol.CredentialCreation) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		b, err := json.Marshal(c)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		if _, err := w.Write(b); err != nil {
			log.Println(err)
			panic(err)
		}
	})
}

func UnmarshalCredentialCreation(v interface{}) (protocol.CredentialCreation, error) {
	return protocol.CredentialCreation{}, fmt.Errorf("%T is not a parsable", v)
}
