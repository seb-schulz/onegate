package model

import (
	"testing"

	"github.com/seb-schulz/onegate/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func TestCreateUser(t *testing.T) {
	tx := openDb().Begin()
	defer tx.Rollback()

	if err := tx.Transaction(CreateUser(nil, nil)); err == nil {
		t.Fatal("failed because error expected")
	}

	user, session := User{}, Session{UserID: &[]uint{1}[0]}
	if err := tx.Transaction(CreateUser(&user, &session)); err == nil {
		t.Fatalf("failed because error expected: user=%v, session=%v", user, session)
	}

	user, session = User{}, Session{}
	if err := tx.Transaction(CreateUser(&user, &session)); err != nil {
		t.Fatalf("failed because error expected: user=%v, session=%v, err=%v", user, session, err)
	}

	if user.ID != *session.UserID {
		t.Fatalf("user ID are not the same: %v != %v", user.ID, session.UserID)
	}

}
