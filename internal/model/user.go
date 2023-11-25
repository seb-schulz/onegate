package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PasskeyID string `gorm:"type:BLOB(16);default:RANDOM_BYTES(16);not null"`
}

const AnonymousUserID = -1
