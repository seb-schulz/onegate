package sessionmgr

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"time"

	"github.com/google/uuid"
)

type (
	Token struct {
		UUID      uuid.UUID
		CreatedAt time.Time
		salt      [4]byte
		sig       func() hash.Hash
	}

	tokenSigner interface {
		sign(key []byte) ([]byte, error)
	}

	tokenParser interface {
		parse(key []byte, token []byte) error
	}

	tokenInitializer interface {
		initialize()
	}

	tokenizer interface {
		tokenInitializer
		tokenSigner
		tokenParser
		fmt.Stringer
	}
)

func newToken() tokenizer {
	return &Token{sig: sha256.New}
}

func (s *Token) initialize() {
	s.UUID = uuid.New()
	s.CreatedAt = time.Now().Truncate(time.Second)

	if _, err := rand.Read(s.salt[:]); err != nil {
		panic(err)
	}
}

func (s *Token) MarshalBinary() ([]byte, error) {
	r, err := s.UUID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	r = binary.BigEndian.AppendUint32(r, uint32(s.CreatedAt.Unix()))
	r = append(s.salt[:], r...)
	return r, nil
}

func (s *Token) UnmarshalBinary(data []byte) error {
	if len(data) > 24 {
		return fmt.Errorf("length of data does not fit")
	}

	u, err := uuid.FromBytes(data[4:20])
	if err != nil {
		return err
	}
	s.UUID = u

	s.CreatedAt = time.Unix(int64(binary.BigEndian.Uint32(data[20:])), 0)

	s.salt = [4]byte(data[:4])

	return nil
}

func (s *Token) sign(key []byte) ([]byte, error) {
	h := hmac.New(s.sig, append(key, s.salt[:]...))

	data, err := s.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}

	if _, err := h.Write(data); err != nil {
		return []byte{}, err
	}

	return append(data, h.Sum(nil)...), nil
}

func (s *Token) size() int {
	x, _ := s.MarshalBinary()
	return len(x)
}

func (s *Token) parse(key []byte, token []byte) error {
	binSize := s.size()
	if len(token) != binSize+s.sig().Size() {
		return fmt.Errorf("token is tampered")
	}

	// Split token into payload part and siganture part
	payload, signature := token[:binSize], token[binSize:]

	h := hmac.New(s.sig, append(key, token[:4]...))
	if _, err := h.Write(payload); err != nil {
		return err
	}

	if !hmac.Equal(signature, h.Sum(nil)) {
		return fmt.Errorf("token is tampered")
	}

	if err := s.UnmarshalBinary(payload); err != nil {
		return err
	}

	return nil
}

func (s *Token) String() string {
	return fmt.Sprint(s.UUID)
}
