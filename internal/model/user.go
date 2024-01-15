package model

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/internal/database"
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

func FirstUser(ctx context.Context) (*User, error) {
	s := Session{ID: sessionmgr.FromContext(ctx).UUID}
	r := database.FromContext(ctx).Preload("User").First(&s)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, r.Error
	}
	return &s.User, nil
}

func CreateUser(ctx context.Context, name string) (*User, error) {
	return database.Transaction(ctx, func(tx *gorm.DB) (*User, error) {
		user := User{Name: name}

		if r := tx.Create(&user); r.Error != nil {
			return nil, r.Error
		}

		token := sessionmgr.FromContext(ctx)
		if r := tx.FirstOrCreate(&Session{
			ID:   token.UUID,
			User: user,
		}); r.Error != nil {
			return nil, r.Error
		}
		return &user, nil
	})
}

type LoginOpt struct {
	UserID     *uint
	Credential *Credential
	Tx         *gorm.DB
}

func (opt *LoginOpt) setUserID(userID *uint) error {
	if opt.UserID != nil {
		return opt.verifyUserID(userID)
	}
	*userID = opt.Credential.UserID
	return nil
}

func (opt *LoginOpt) verifyUserID(userID *uint) error {
	r := opt.Tx.Model(&User{}).Where("id = ?", *opt.UserID).Limit(1).Pluck("id", userID)
	if r.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func LoginUser(ctx context.Context, opt LoginOpt) error {
	if opt.UserID == nil && opt.Credential == nil {
		return fmt.Errorf("either UserID or Credetial must be defined")
	}

	if opt.UserID != nil && opt.Credential != nil {
		return fmt.Errorf("UserID and Credential are mutually exclusive")
	}

	if opt.Tx == nil {
		opt.Tx = database.FromContext(ctx)
	}

	if err := opt.Tx.Transaction(func(tx *gorm.DB) error {
		var userID uint
		if err := opt.setUserID(&userID); err != nil {
			return err
		}

		if r := tx.FirstOrCreate(&Session{
			ID:     sessionmgr.FromContext(ctx).UUID,
			UserID: userID,
		}); r.Error != nil {
			return r.Error
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
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
