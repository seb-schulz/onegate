package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db_dsn = os.Getenv("DB_DSN")

func init() {
	if db_dsn == "" {
		db_dsn = "onegate:.test.@tcp(db:3306)/onegate?charset=utf8&parseTime=True"
	}
	// fmt.Println(os.Environ())
	fmt.Println(db_dsn)
	// viper.ReadConfig(bytes.NewBuffer([]byte("sessionKey: LnRlc3Qu")))
}

func openDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(db_dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
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
