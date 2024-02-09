package auth

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"

	"go.pact.im/x/option"
	"go.pact.im/x/phcformat"
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

	hash := newPBKDF2Key([]byte("abc"), 1024, 32, sha1.New)

	cID, cs, err := createClient(database.WithContext(context.Background(), tx), hash, "http://localhost:9000/cb")
	if err != nil {
		t.Errorf("cannot create client: %v", err)
	}

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

	if err := c.VerifyClientSecret(cs); err != nil {
		t.Errorf("cannot verify secret: %v", err)
	}

	// t.Errorf("client=%#v, secret=%s", c, c.ClientSecret)
}

func TestClientSecretKeyer(t *testing.T) {
	h := newPBKDF2Key([]byte{1, 2, 3}, 1, 32, sha1.New)

	key := []byte("a")
	log.Println(h.phcString(key))
	fakeClient := Client{ClientSecret: h.phcString(key)}

	if err := fakeClient.VerifyClientSecret(base64.URLEncoding.EncodeToString(key)); err != nil {
		t.Errorf("verification failed with key=%s and err=%v", key, err)
	}

	// if clientSecretHasherType(hashedKey[0]) != clientSecretHasherPBKDF2 {
	// 	t.Errorf("first byte is not expected type")
	// }
	// t.Error(fakeClient)
}

func TestPHCFormat(t *testing.T) {
	h := phcformat.MustParse("$name$v=42$k=v$salt$hash")
	fmt.Println(h)
	// fmt.Println(h.ID)
	// fmt.Println(option.UnwrapOrZero(h.Version))
	// fmt.Println(option.UnwrapOrZero(h.Params))
	// fmt.Println(option.UnwrapOrZero(h.Salt))
	// fmt.Println(option.UnwrapOrZero(h.Output))

	h = phcformat.MustParse("$pbkdf2-sha1$i=1,k=32$AQID$gbkoCyl+kgPHGtDI1sbgzYTjh1vawpwGBu7TsfyyQ7Y")
	fmt.Println(h)
	fmt.Println(h.ID)
	fmt.Println(option.UnwrapOrZero(h.Salt))
	fmt.Println(option.UnwrapOrZero(h.Output))
	fmt.Println(option.UnwrapOrZero(h.Params))

	for it := phcformat.IterParams(option.UnwrapOrZero(h.Params)); it.Valid; it = it.Next() {
		fmt.Println(it.Name, it.Value)
	}
}
