package model

import (
	"encoding/base64"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PasskeyID string `gorm:"type:BLOB(16);default:RANDOM_BYTES(16);not null"`
}

func (u User) Base64PasskeyID() string {
	return base64.StdEncoding.EncodeToString([]byte(u.PasskeyID))
}
