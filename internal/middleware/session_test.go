package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seb-schulz/onegate/internal/utils"
	"gorm.io/gorm"
)

func openDb() *gorm.DB {
	db, err := utils.OpenDatabase()
	if err != nil {
		panic(err)
	}
	return db
}

func checkCookie(resp *http.Response, name string) bool {
	for _, c := range resp.Cookies() {
		if c.Name == name {
			return true
		}
	}
	return false
}

func TestSessionMiddleware(t *testing.T) {
	tx := openDb().Begin()
	defer tx.Rollback()

	handler := SessionMiddleware(tx)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !checkCookie(w.Result(), "session") {
		t.FailNow()
	}

	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.AddCookie(w.Result().Cookies()[0])
	handler.ServeHTTP(w, req)
	if !checkCookie(w.Result(), "session") {
		t.FailNow()
	}
}
