package auth

import (
	"context"
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"
)

func TestClientByClientID(t *testing.T) {
	db, err := database.Open()
	if err != nil {
		panic(err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	id, _ := uuid.NewUUID()

	tx.Create(&Client{ID: id})

	c, err := clientByClientID(database.WithContext(context.Background(), tx), fmt.Sprint(id))
	if err != nil {
		t.Errorf("cannot get client: %v", err)
	}

	if c.ClientID() != fmt.Sprint(id) {
		t.Errorf("got client id %#v instead of %#v", c.ClientID(), id)
	}

	// t.Errorf("client: %#v", c)
}

func TestCreateClient(t *testing.T) {
	db, err := database.Open()
	if err != nil {
		panic(err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	hash := newPBKDF2Key([]byte("a"), 1024, 32, sha1.New)

	cID, cs, err := createClient(database.WithContext(context.Background(), tx), hash, "http://localhost:9000/cb")
	if cID == "" {
		t.Errorf("client ID was empty string")
	}

	if cs == "" {
		t.Errorf("client secret was empty string")
	}

	c := Client{}
	tx.First(&c, "id = ?", cID)
	if c.RedirectURI() != "http://localhost:9000/cb" {
		t.Errorf("redirect URI not matching")
	}

	if err := c.verifyClientSecret(hash, cs); err != nil {
		t.Errorf("cannot verify secret: %v", err)
	}

	// t.Errorf("client=%#v, secret=%s", c, base64.URLEncoding.EncodeToString(c.ClientSecret))
}

func TestClientSecretKeyer(t *testing.T) {
	h := newPBKDF2Key([]byte{1}, 1, 32, sha1.New)

	key := []byte("a")
	if hashedKey := h.Key(key); !h.Verify(hashedKey, key) {
		t.Errorf("verification failed: %s", key)
	}

}
