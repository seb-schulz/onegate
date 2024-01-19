package graph

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB                      *gorm.DB
	WebAuthn                *webauthn.WebAuthn
	UserMgr                 *sessionmgr.StorageManager[*model.User]
	UserRegistrationEnabled bool
}
