package model

import "gorm.io/gorm"

type UserJwtConverter interface {
	Subject() string
}

type User struct {
	gorm.Model
	PasskeyID string `gorm:"type:BLOB(16);default:RANDOM_BYTES(16);not null"`
}

type anonymousUser struct{}

func (u anonymousUser) Subject() string {
	return "anon"
}

var AnonymousUser = anonymousUser{}
