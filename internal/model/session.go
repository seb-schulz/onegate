package model

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"math/rand"
	"strings"
	"time"

	"github.com/seb-schulz/onegate/internal/config"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID *uint
	User   *User
}

func nonce() []byte {
	characters := "ABCDEFGHIJKLMOPQRSTUVWXYZabcdefghijklmopqrstuvwxyz0123456789"
	id_runes := make([]byte, 4)
	for i := range id_runes {
		id_runes[i] = characters[rand.Intn(len(characters))]
	}
	return id_runes
}

func newHMAC() hash.Hash {
	return hmac.New(sha256.New, []byte(config.Default.Session.Key))
}

func generateToken(id uint, nonce []byte) string {
	mac := newHMAC()
	idStr := fmt.Sprintf("%x", id)
	mac.Write([]byte(idStr))
	mac.Write(nonce)

	return fmt.Sprintf("%s-%s-%s", idStr, nonce, hex.EncodeToString(mac.Sum(nil)))

}

func (s Session) Token() string {
	return generateToken(s.ID, nonce())
}

func (s Session) String() string {
	uID := "nil"
	if s.UserID != nil {
		uID = fmt.Sprintf("%v", *s.UserID)
	}

	u := "nil"
	if s.User != nil {
		u = fmt.Sprintf("User{ID: %v}", s.User.ID)
	}
	return fmt.Sprintf("Session{ID: %d, UserID: %v, User: %v)", s.ID, uID, u)
}

func (s Session) IsActive() bool {
	return time.Since(s.UpdatedAt) <= config.Default.Session.ActiveFor
}

func getSessionIDByToken(token string) (uint, error) {
	slicedToken := strings.Split(token, "-")
	if len(slicedToken) != 3 {
		return 0, fmt.Errorf("invalid token")
	}

	mac := newHMAC()
	mac.Write([]byte(slicedToken[0]))
	mac.Write([]byte(slicedToken[1]))

	tokenSig, err := hex.DecodeString(slicedToken[2])
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	if !hmac.Equal(tokenSig, mac.Sum(nil)) {
		return 0, fmt.Errorf("invalid token")
	}

	var id uint
	if _, err := fmt.Sscanf(slicedToken[0], "%x", &id); err != nil {
		return 0, err
	}

	return id, nil
}

func FirstSessionByToken(db *gorm.DB, token string, session *Session) error {
	id, err := getSessionIDByToken(token)
	if err != nil {
		return err
	}

	db.Preload("User").FirstOrCreate(session, "id = ?", id)
	return nil
}

func CreateSession(db *gorm.DB, session *Session) {
	db.Create(session)
}

func AllSessionByUserID(db *gorm.DB, userID uint) ([]*Session, error) {
	sessions := []*Session{}
	r := db.Where("user_id = ?", userID).Find(&sessions)

	return sessions, r.Error
}

func DeleteSessionByUserID(userID uint, id string) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		s := Session{}
		if result := tx.Where("user_id = ?", userID).First(&s, id); result.Error != nil {
			return result.Error
		}
		tx.Delete(&s)
		return nil
	}
}
