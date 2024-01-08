package server

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash"
	"time"

	"github.com/google/uuid"
)

type sessionToken struct {
	UUID      uuid.UUID
	CreatedAt time.Time
	Salt      [4]byte
}

var sessionBinarySize = 0

func init() {
	x, _ := (&sessionToken{}).MarshalBinary()
	sessionBinarySize = len(x)
}

func newSessionToken() *sessionToken {
	s := sessionToken{
		UUID:      uuid.Must(uuid.NewRandom()),
		CreatedAt: time.Now().Truncate(time.Second),
	}

	if _, err := rand.Read(s.Salt[:]); err != nil {
		panic(err)
	}

	return &s
}

func (s *sessionToken) MarshalBinary() ([]byte, error) {
	r, err := s.UUID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	r = binary.BigEndian.AppendUint32(r, uint32(s.CreatedAt.Unix()))
	r = append(s.Salt[:], r...)
	return r, nil
}

func (s *sessionToken) UnmarshalBinary(data []byte) error {
	if len(data) > 24 {
		return fmt.Errorf("length of data does not fit")
	}

	u, err := uuid.FromBytes(data[4:20])
	if err != nil {
		return err
	}
	s.UUID = u

	s.CreatedAt = time.Unix(int64(binary.BigEndian.Uint32(data[20:])), 0)

	s.Salt = [4]byte(data[:4])

	return nil
}

func (s *sessionToken) signedToken(key []byte, sig func() hash.Hash) (string, error) {
	h := hmac.New(sig, append(key, s.Salt[:]...))

	data, err := s.MarshalBinary()
	if err != nil {
		return "", err
	}

	if _, err := h.Write(data); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(append(data, h.Sum(nil)...)), nil
}

func parseToken(key []byte, sig func() hash.Hash, token string) (*sessionToken, error) {
	rawToken, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	if len(rawToken) != sessionBinarySize+sig().Size() {
		return nil, fmt.Errorf("token is tampered")
	}

	// Split token into payload part and siganture part
	payload, signature := rawToken[:sessionBinarySize], rawToken[sessionBinarySize:]

	h := hmac.New(sig, append(key, rawToken[:4]...))
	if _, err := h.Write(payload); err != nil {
		return nil, err
	}

	if !hmac.Equal(signature, h.Sum(nil)) {
		return nil, fmt.Errorf("token is tampered")
	}

	s := sessionToken{}
	if err := s.UnmarshalBinary(payload); err != nil {
		return nil, err
	}

	return &s, nil
}
