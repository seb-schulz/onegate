package utils

import (
	"fmt"

	"github.com/seb-schulz/onegate/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type optFuncs func(*gorm.DB) *gorm.DB

func OpenDatabase(opts ...optFuncs) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Config.DB.Dsn), &gorm.Config{})

	if err != nil {
		return db, fmt.Errorf("failed to connect to database: %v", err)
	}

	for _, fn := range opts {
		db = fn(db)
	}

	return db, nil
}

func WithDebugOption(debug bool) optFuncs {
	return func(tx *gorm.DB) *gorm.DB {
		if debug {
			return tx.Debug()
		}
		return tx
	}
}
