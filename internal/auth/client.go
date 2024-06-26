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
	"golang.org/x/crypto/argon2"
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
	ClientID() uuid.UUID
	ClientSecretVerifier
	redirecter
}

type clientByClientIDFn func(ctx context.Context, clientID string) (client, error)

type Client struct {
	ID                  uuid.UUID `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	Description         string         `gorm:"type:VARCHAR(255);not null"`
	ClientSecret        string         `gorm:"type:VARCHAR(255);not null"`
	InternalRedirectURI string         `gorm:"column:redirect_uri;type:VARCHAR(255);not null"`
}

func (c *Client) ClientID() uuid.UUID {
	return c.ID
}

func (c *Client) RedirectURI() string {
	return fmt.Sprint(c.InternalRedirectURI)
}

func (c *Client) VerifyClientSecret(s string) error {
	decodedSecret, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}

	var hasher ClientSecretHasher

	phcHash := phcformat.MustParse(c.ClientSecret)
	switch phcHash.ID {
	case "pbkdf2-sha1":
		hasher = newPBKDF2KeyFromPHCWithSha1(phcHash)
	case "argon2id":
		hasher = newArgon2IdFromPHC(phcHash)
	default:
		return fmt.Errorf("hash algorithm unknown")
	}

	if subtle.ConstantTimeCompare([]byte(option.UnwrapOrZero(phcHash.Output)), hasher.Key(decodedSecret)) != 1 {
		log.Printf("%v != %s", option.UnwrapOrZero(phcHash.Output), s)
		return fmt.Errorf("verification failed")
	}

	return nil
}

func clientByClientID(ctx context.Context, clientID string) (client, error) {
	var c Client
	r := database.FromContext(ctx).First(&c, "id = ?", clientID)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, r.Error
	}
	return &c, nil
}

func CreateClient(ctx context.Context, clientSecretHash ClientSecretHasher, desc, redirectURL string) (clientID string, clientSecret string, err error) {
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
		Description:         desc,
	}

	r := database.FromContext(ctx).Create(&client)
	if r.Error != nil {
		return "", "", r.Error
	}

	return fmt.Sprint(client.ClientID()), base64.URLEncoding.EncodeToString(randSecret[:]), nil
}

type pbkdf2Key struct {
	salt []byte
	iter int
	hash func() hash.Hash
}

func newPBKDF2Key(salt []byte, iter int, h func() hash.Hash) *pbkdf2Key {
	return &pbkdf2Key{salt, iter, h}
}

func newPBKDF2KeyFromPHCWithSha1(hash phcformat.Hash) *pbkdf2Key {
	var iter int
	for it := phcformat.IterParams(option.UnwrapOrZero(hash.Params)); it.Valid; it = it.Next() {
		var err error

		switch it.Name {
		case "i":
			iter, err = strconv.Atoi(it.Value)
			if err != nil {
				panic(fmt.Errorf("invalid iter format: %v", err))
			}
		default:
			panic("invalid format")
		}
	}
	rawSalt, err := base64.RawStdEncoding.DecodeString(option.UnwrapOrZero(hash.Salt))
	if err != nil {
		panic("invalid salt")
	}
	return newPBKDF2Key(rawSalt, iter, sha1.New)
}

func (h *pbkdf2Key) rawKey(bKey []byte) []byte {
	return pbkdf2.Key(bKey, h.salt, h.iter, 30, sha1.New)
}

func (h *pbkdf2Key) Key(bKey []byte) []byte {
	s := base64.RawStdEncoding.EncodeToString(h.rawKey(bKey))
	return []byte(s)
}

func (h *pbkdf2Key) phcString(key []byte) string {
	return fmt.Sprintf("$pbkdf2-sha1$i=%d$%s$%s", h.iter, base64.RawStdEncoding.EncodeToString(h.salt), base64.RawStdEncoding.EncodeToString(h.rawKey(key)))
}

type argon2IdKey struct {
	salt    []byte
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func newArgon2Id(salt []byte, time, memory uint32, threads uint8) *argon2IdKey {
	return &argon2IdKey{salt, time, memory, threads, 30}
}

func newArgon2IdFromPHC(hash phcformat.Hash) *argon2IdKey {
	var threads, time, memory uint64

	for it := phcformat.IterParams(option.UnwrapOrZero(hash.Params)); it.Valid; it = it.Next() {
		var err error

		switch it.Name {
		case "m":
			memory, err = strconv.ParseUint(it.Value, 10, 32)
			if err != nil {
				panic(fmt.Errorf("invalid iter format: %v", err))
			}
		case "t":
			time, err = strconv.ParseUint(it.Value, 10, 32)
			if err != nil {
				panic(fmt.Errorf("invalid iter format: %v", err))
			}
		case "p":
			threads, err = strconv.ParseUint(it.Value, 10, 32)
			if err != nil {
				panic(fmt.Errorf("invalid iter format: %v", err))
			}
		default:
			panic("invalid format")
		}
	}
	rawSalt, err := base64.RawStdEncoding.DecodeString(option.UnwrapOrZero(hash.Salt))
	if err != nil {
		panic("invalid salt")
	}
	return newArgon2Id(rawSalt, uint32(time), uint32(memory), uint8(threads))
}

func (h *argon2IdKey) rawKey(bKey []byte) []byte {
	return argon2.IDKey(bKey, h.salt, h.time, h.memory, h.threads, h.keyLen)
}

func (h *argon2IdKey) Key(bKey []byte) []byte {
	s := base64.RawStdEncoding.EncodeToString(h.rawKey(bKey))
	return []byte(s)
}

func (h *argon2IdKey) phcString(key []byte) string {
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.memory, h.time, h.threads, base64.RawStdEncoding.EncodeToString(h.salt), base64.RawStdEncoding.EncodeToString(h.rawKey(key)))
}

var readRand = func(b []byte) error {
	_, err := rand.Read(b)
	return err
}

func NewClientSecretHasher() ClientSecretHasher {
	randSalt := make([]byte, 15)
	if err := readRand(randSalt); err != nil {
		panic("cannot generate salt")

	}
	return newArgon2Id(randSalt, 2, 64*1024, 4)
}
