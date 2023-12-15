package graph

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

//go:generate go run github.com/99designs/gqlgen generate
//go:generate /bin/bash -c "(cd $(pwd)/../internal/ui/_client/ && npm run compile)"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB       *gorm.DB
	WebAuthn *webauthn.WebAuthn
}
