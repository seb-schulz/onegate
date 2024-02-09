package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"
	"go.pact.im/x/option"
	"go.pact.im/x/phcformat"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

type redirecter interface {
	RedirectURI() string
}

type ClientSecretHasher interface {
	Key([]byte) []byte
	phcString([]byte) string
}

type ClientSecretVerifier interface {
	VerifyClientSecret(string) error
}

type client interface {
	ClientID() string
	ClientSecretVerifier
	redirecter
}

type clientByClientIDFn func(ctx context.Context, clientID string) (client, error)

type Client struct {
	ID                  uuid.UUID `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	ClientSecret        string
	InternalRedirectURI string `gorm:"column:redirect_uri"`
}

func (c *Client) ClientID() string {
	return fmt.Sprint(c.ID)
}

func (c *Client) RedirectURI() string {
	return fmt.Sprint(c.InternalRedirectURI)
}

func (c *Client) VerifyClientSecret(s string) error {
	decodedSecret, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}

	phcHash := phcformat.MustParse(c.ClientSecret)
	switch phcHash.ID {
	case "pbkdf2-sha1":
		hasher := newPBKDF2KeyFromPHCWithSha1(phcHash)
		if subtle.ConstantTimeCompare([]byte(option.UnwrapOrZero(phcHash.Output)), hasher.Key(decodedSecret)) != 1 {
			log.Printf("%v != %s", option.UnwrapOrZero(phcHash.Output), s)
			return fmt.Errorf("verification failed")
		}
	default:
		return fmt.Errorf("hash algorithm unknown")
	}
	return nil
}

func clientByClientID(ctx context.Context, clientID string) (client, error) {
	var c Client
	r := database.FromContext(ctx).First(&c)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, r.Error
	}
	return &c, nil
}

func createClient(ctx context.Context, clientSecretHash ClientSecretHasher, redirectURL string) (clientID string, clientSecret string, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("cannot generate uuid: %v", err))
	}

	// TODO: Provide stable salt value
	randSecret := make([]byte, 32)
	if _, err := rand.Read(randSecret); err != nil {
		return "", "", err
	}

	client := Client{
		ID:                  id,
		ClientSecret:        clientSecretHash.phcString(randSecret[:]),
		InternalRedirectURI: redirectURL,
	}

	r := database.FromContext(ctx).Create(&client)
	if r.Error != nil {
		return "", "", r.Error
	}

	return client.ClientID(), base64.URLEncoding.EncodeToString(randSecret[:]), nil
}

type pbkdf2Key struct {
	salt   []byte
	iter   int
	keyLen int
	hash   func() hash.Hash
}

func newPBKDF2Key(salt []byte, iter, keyLen int, h func() hash.Hash) *pbkdf2Key {
	return &pbkdf2Key{salt, iter, keyLen, h}
}

func newPBKDF2KeyFromPHCWithSha1(hash phcformat.Hash) *pbkdf2Key {
	var iter, keyLen int
	for it := phcformat.IterParams(option.UnwrapOrZero(hash.Params)); it.Valid; it = it.Next() {
		var err error

		switch it.Name {
		case "i":
			iter, err = strconv.Atoi(it.Value)
			if err != nil {
				panic(fmt.Errorf("invalid iter format: %v", err))
			}
		case "k":
			keyLen, err = strconv.Atoi(it.Value)
			if err != nil {
				panic(fmt.Errorf("invalid keyLen format: %v", err))
			}
		default:
			panic("invalid format")
		}
	}
	rawSalt, err := base64.RawStdEncoding.DecodeString(option.UnwrapOrZero(hash.Salt))
	if err != nil {
		panic("invalid salt")
	}
	return newPBKDF2Key(rawSalt, iter, keyLen, sha1.New)
}

func (h *pbkdf2Key) rawKey(bKey []byte) []byte {
	return pbkdf2.Key(bKey, h.salt, h.iter, h.keyLen, sha1.New)
}

func (h *pbkdf2Key) Key(bKey []byte) []byte {
	s := base64.RawStdEncoding.EncodeToString(h.rawKey(bKey))
	return []byte(s)
}

func (h *pbkdf2Key) phcString(key []byte) string {
	return fmt.Sprintf("$pbkdf2-sha1$i=%d,k=%d$%s$%s", h.iter, h.keyLen, base64.URLEncoding.EncodeToString(h.salt), base64.RawStdEncoding.EncodeToString(h.rawKey(key)))
}
