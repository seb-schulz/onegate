package server

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash"
	"net/http"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
)

type sessionToken struct {
	UUID      uuid.UUID
	CreatedAt time.Time
	Salt      [4]byte
	sig       func() hash.Hash
}

type sessionTokenSigner interface {
	sign(key []byte) ([]byte, error)
}

type sessionTokenParser interface {
	parse(key []byte, token []byte) error
}

type sessionTokenInitializer interface {
	initialize()
}

type sessionTokenizer interface {
	sessionTokenInitializer
	sessionTokenSigner
	sessionTokenParser
}

var sessionBinarySize = 0

func init() {
	x, _ := (&sessionToken{}).MarshalBinary()
	sessionBinarySize = len(x)
}

func newSessionToken() sessionTokenizer {
	return &sessionToken{sig: sha256.New}
}

func (s *sessionToken) initialize() {
	s.UUID = uuid.Must(uuid.NewRandom())
	s.CreatedAt = time.Now().Truncate(time.Second)

	if _, err := rand.Read(s.Salt[:]); err != nil {
		panic(err)
	}
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

func (s *sessionToken) sign(key []byte) ([]byte, error) {
	h := hmac.New(s.sig, append(key, s.Salt[:]...))

	data, err := s.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}

	if _, err := h.Write(data); err != nil {
		return []byte{}, err
	}

	return append(data, h.Sum(nil)...), nil
}

func (s *sessionToken) parse(key []byte, token []byte) error {
	if len(token) != sessionBinarySize+s.sig().Size() {
		return fmt.Errorf("token is tampered")
	}

	// Split token into payload part and siganture part
	payload, signature := token[:sessionBinarySize], token[sessionBinarySize:]

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

type contextSessionKeyType struct{ string }

var contextSessionToken = contextSessionKeyType{"session"}

type sessionMiddleware struct {
	key      []byte
	newToken func() sessionTokenizer
}

func (s *sessionMiddleware) setCookie(w http.ResponseWriter, token sessionTokenSigner) {
	sToken, err := token.sign(s.key)
	if err != nil {
		panic(err) // signing token should not fail
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    base64.RawURLEncoding.EncodeToString(sToken),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (s *sessionMiddleware) tokenFromCookie(req *http.Request, token sessionTokenParser) error {
	cookie, err := req.Cookie("session")
	if err != nil {
		return err
	}

	rawToken, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return err
	}

	if err := token.parse(s.key, rawToken); err != nil {
		logger := httplog.LogEntry(req.Context())
		logger.Warn(fmt.Sprintf("cannot parse raw token: %v", err))
		return err
	}
	return nil
}

func (s *sessionMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := s.newToken()

		if err := s.tokenFromCookie(r, token); err != nil {
			token.initialize()
			s.setCookie(w, token)
		}

		ctx := context.WithValue(r.Context(), contextSessionToken, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func mustSessionTokenFromContext(ctx context.Context) sessionTokenizer {
	raw, ok := ctx.Value(contextSessionToken).(sessionTokenizer)
	if !ok {
		panic("session token does not exist")
	}
	return raw
}
