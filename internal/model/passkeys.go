package model

import "gorm.io/gorm"

type Passkeys struct {
	gorm.Model
	UserID        int
	User          User
	Username      string `gorm:"type:VARCHAR(255);not null"`
	PublicKeySpki []byte `gorm:"type:BLOB"`
	Backup        bool
}
