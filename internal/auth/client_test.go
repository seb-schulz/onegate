package auth

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"
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

	hash := newPBKDF2Key([]byte("abc"), 1024, sha1.New)

	cID, cs, err := CreateClient(database.WithContext(context.Background(), tx), hash, "hello world", "http://localhost:9000/cb")
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

func FuzzClientSecretKeyer(f *testing.F) {

	for i := 0; i < 100; i++ {
		f.Add(rand.Int())
	}
	f.Fuzz(func(t *testing.T, seed int) {
		gen := rand.New(rand.NewSource(int64(seed)))

		pbkdf2Seed := make([]byte, 3)
		gen.Read(pbkdf2Seed)
		for _, h := range []ClientSecretHasher{
			newPBKDF2Key(pbkdf2Seed, 1, sha1.New),
			newArgon2Id(pbkdf2Seed, 1, 1, 1),
		} {
			key := make([]byte, 32)
			gen.Read(key)

			fakeClient := Client{ClientSecret: h.phcString(key)}

			if err := fakeClient.VerifyClientSecret(base64.URLEncoding.EncodeToString(key)); err != nil {
				t.Errorf("verification failed with key=%s and err=%v", key, err)
				t.FailNow()
			}
			// t.Error(fakeClient)
		}
	})

}

func TestNewClientSecretHasher(t *testing.T) {
	gen := rand.New(rand.NewSource(int64(1)))

	_orig := readRand
	readRand = func(b []byte) error {
		_, err := gen.Read(b)
		return err
	}

	defer func() {
		readRand = _orig
	}()

	for i := 1; i <= 10; i++ {
		key := make([]byte, 32)
		gen.Read(key)

		phcHash := NewClientSecretHasher().phcString(key)
		t.Logf("%v", phcHash)
		fakeClient := Client{ClientSecret: phcHash}

		if err := fakeClient.VerifyClientSecret(base64.URLEncoding.EncodeToString(key)); err != nil {
			t.Errorf("verification failed with key=%s and err=%v", key, err)
		}
		// t.Error(fakeClient)
	}

}

func TestMustParsePhcformat(t *testing.T) {
	t.Skip()
	phcformat.MustParse("$argon2id$v=19$m=65536,t=1,p=4,k=32$Ii1+3omeYiOsWbIJ2/plPgG9$8Gw1SSuNrdPsCzkH9O+eXBIsOomJJ0zIdq5G5EaGtIE")
}
