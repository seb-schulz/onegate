package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
)

func TestAuthorizationMgrCreate(t *testing.T) {
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

	authMgr := authorizationMgr{
		StorageManager: sessionmgr.NewStorage("authorization", FirstAuthorization),
	}

	if err := authMgr.create(ctx, &client, "state", "CodeChallenge"); err != nil {
		t.Errorf("failed to create authorization: %v", err)
	}

	route := chi.NewRouter()
	route.Use(authMgr.Handler)
	route.Get("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Ok")

		authReq := authMgr.FromContext(r.Context())
		if authReq == nil {
			t.Errorf("cannot find authorization from database")
		}
	}))

	w := httptest.NewRecorder()
	route.ServeHTTP(w, httptest.NewRequest("GET", "/foo", nil).WithContext(ctx))

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.FailNow()
	}
}
