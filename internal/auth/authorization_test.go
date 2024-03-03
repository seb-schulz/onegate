package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
	"gorm.io/gorm"
)

func TestCreateAuthorization(t *testing.T) {
	db, err := database.Open()
	if err != nil {
		panic(err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	sessionToken := sessionmgr.Token{UUID: uuid.New()}
	ctx := sessionmgr.ToContext(database.WithContext(context.Background(), tx), &sessionToken)

	client := Client{
		ID: uuid.New(),
	}
	tx.FirstOrCreate(&client)

	if err := createAuthorization(ctx, &client, "state", "CodeChallenge"); err != nil {
		t.Errorf("failed to create authorization: %v", err)
	}

	// TODO: assert entries after authorization creation
	// Extract code below into login handler test

	route := chi.NewRouter()
	route.Get("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Ok")

		_, err := FirstAuthorization(r.Context())
		if err != nil {
			t.Errorf("cannot find authorization: %v", err)
		}
	}))

	w := httptest.NewRecorder()
	route.ServeHTTP(w, httptest.NewRequest("GET", "/foo", nil).WithContext(ctx))

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.FailNow()
	}
}

func TestAuthorizationByCode(t *testing.T) {
	gen := rand.New(rand.NewSource(int64(1)))

	readRand := func(size uint) []byte {
		b := make([]byte, size)
		gen.Read(b)
		return b
	}

	db, err := database.Open()
	if err != nil {
		panic(err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	ctx := database.WithContext(context.Background(), tx)

	client := Client{
		ID: uuid.New(),
	}
	tx.FirstOrCreate(&client)

	user := model.User{}
	tx.FirstOrCreate(&user)

	code := readRand(16)

	r := tx.FirstOrCreate(&Authorization{Client: client, InternalCode: code})
	if r.Error != nil {
		t.Errorf("cannot create authorization: %v", r.Error)
	}

	_, err = authorizationByCode(ctx, base64.URLEncoding.EncodeToString([]byte("non existing error")))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("failed but with unexpected error msg: %v", err)
	} else if err == nil {
		t.Errorf("expected error but found authorization")
	}

	_, err = authorizationByCode(ctx, "expected decoding error")
	if err == nil {
		t.Error("expected decoding error")
	}

	fetchedAuth, err := authorizationByCode(ctx, base64.URLEncoding.EncodeToString(code))
	if err != nil {
		t.Errorf("failed to get authorization by code: %v", err)
	}

	t.Logf("fetchedAuth: %v", fetchedAuth)
}
