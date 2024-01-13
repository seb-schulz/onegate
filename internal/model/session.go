package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/config"
	"gorm.io/gorm"
)

type Session struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint
	User      User
}

func (s Session) String() string {
	return fmt.Sprintf("Session{ID: %d, UserID: %v, User: %v)", s.ID, s.UserID, s.User)
}

func (s Session) IsActive() bool {
	return time.Since(s.UpdatedAt) <= config.Config.Session.ActiveFor
}

func DeleteSessionByUserID(userID uint, id uuid.UUID) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		s := Session{ID: id}
		if result := tx.Where("user_id = ?", userID).First(&s); result.Error != nil {
			return result.Error
		}
		tx.Delete(&s)
		return nil
	}
}

func AllSessionByUserID(tx *gorm.DB, userID uint) ([]*Session, error) {
	sessions := []*Session{}
	r := tx.Where("user_id = ?", userID).Find(&sessions)
	return sessions, r.Error
}
