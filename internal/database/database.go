package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/seb-schulz/onegate/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type contextDatabaseKeyType struct{ string }

type optFuncs func(*gorm.DB) *gorm.DB

var ctxDatabaseKey = contextDatabaseKeyType{"DB"}

func Open(opts ...optFuncs) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Config.DB.Dsn), &gorm.Config{})

	if err != nil {
		return db, fmt.Errorf("failed to connect to database: %v", err)
	}

	for _, fn := range opts {
		db = fn(db)
	}

	return db, nil
}

func WithDebug(debug bool) optFuncs {
	return func(tx *gorm.DB) *gorm.DB {
		if debug {
			return tx.Debug()
		}
		return tx
	}
}

func Middleware(tx *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ctxDatabaseKey, tx.WithContext(context.Background()))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromContext(ctx context.Context) *gorm.DB {
	raw, ok := ctx.Value(ctxDatabaseKey).(*gorm.DB)
	if !ok {
		panic("database connection does not exist on context")
	}
	return raw
}

type transactionOpts struct {
	tx *gorm.DB
}

type TransactionOptFunc func(*transactionOpts)

func WithNestedTransaction(tx *gorm.DB) TransactionOptFunc {
	return func(to *transactionOpts) {
		to.tx = tx
	}
}

func Transaction[T any](ctx context.Context, fn func(*gorm.DB) (T, error), opts ...TransactionOptFunc) (T, error) {
	var result T
	transactionOpts := transactionOpts{FromContext(ctx)}

	for _, opt := range opts {
		opt(&transactionOpts)
	}

	err := transactionOpts.tx.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = fn(tx)
		return err
	})
	return result, err
}
