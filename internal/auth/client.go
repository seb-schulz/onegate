package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"time"

	"github.com/google/uuid"
	"github.com/seb-schulz/onegate/internal/database"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

type redirecter interface {
	RedirectURI() string
}

type client interface {
	ClientID() string
	verifyClientSecret(ClientSecretVerifier, string) error
	redirecter
}

type ClientSecretHasher interface {
	Key([]byte) []byte
}

type ClientSecretVerifier interface {
	Verify([]byte, []byte) bool
}

type clientByClientIDFn func(ctx context.Context, clientID string) (client, error)

type Client struct {
	ID                  uuid.UUID `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	ClientSecret        []byte
	InternalRedirectURI string `gorm:"column:redirect_uri"`
}

func (c *Client) ClientID() string {
	return fmt.Sprint(c.ID)
}

func (c *Client) RedirectURI() string {
	return fmt.Sprint(c.InternalRedirectURI)
}

func (c *Client) verifyClientSecret(clientSecretHash ClientSecretVerifier, s string) error {
	rawSecret, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}

	if !clientSecretHash.Verify(c.ClientSecret, []byte(rawSecret)) {
		return fmt.Errorf("client secret miss match")
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
		ClientSecret:        clientSecretHash.Key(randSecret[:]),
		InternalRedirectURI: redirectURL,
	}

	r := database.FromContext(ctx).Create(&client)
	if r.Error != nil {
		return "", "", r.Error
	}

	return client.ClientID(), base64.URLEncoding.EncodeToString(randSecret[:]), nil
}

type PBKDF2Key struct {
	salt   []byte
	iter   int
	keyLen int
	h      func() hash.Hash
}

func newPBKDF2Key(salt []byte, iter, keyLen int, h func() hash.Hash) *PBKDF2Key {
	return &PBKDF2Key{salt, iter, keyLen, h}
}

func (h *PBKDF2Key) Key(key []byte) []byte {
	return pbkdf2.Key(key, h.salt, h.iter, h.keyLen, h.h)
}

func (h *PBKDF2Key) Verify(orig, key []byte) bool {
	k := pbkdf2.Key(key, h.salt, h.iter, h.keyLen, h.h)
	return subtle.ConstantTimeCompare(orig, k) == 1
}
