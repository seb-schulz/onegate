package server

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type sessionToken struct {
	UUID      uuid.UUID
	CreatedAt time.Time
	Salt      [4]byte
}

type sessionTokenSigner interface {
	sign(key []byte, sig func() hash.Hash) ([]byte, error)
}

type sessionTokenParser interface {
	parse(key []byte, sig func() hash.Hash, token []byte) error
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

func (s *sessionToken) sign(key []byte, sig func() hash.Hash) ([]byte, error) {
	h := hmac.New(sig, append(key, s.Salt[:]...))

	data, err := s.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}

	if _, err := h.Write(data); err != nil {
		return []byte{}, err
	}

	return append(data, h.Sum(nil)...), nil
}

func (s *sessionToken) parse(key []byte, sig func() hash.Hash, token []byte) error {
	if len(token) != sessionBinarySize+sig().Size() {
		return fmt.Errorf("token is tampered")
	}

	// Split token into payload part and siganture part
	payload, signature := token[:sessionBinarySize], token[sessionBinarySize:]

	h := hmac.New(sig, append(key, token[:4]...))
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
	signer   func() hash.Hash
	newToken func() sessionTokenizer
}

func (s *sessionMiddleware) setCookie(w http.ResponseWriter, token sessionTokenSigner) {
	sToken, err := token.sign(s.key, s.signer)
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

	return token.parse(s.key, s.signer, rawToken)
}

func (s *sessionMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := s.newToken()

		if err := s.tokenFromCookie(r, token); err != nil {
			token.initialize()
			s.setCookie(w, token)
			ctx := context.WithValue(r.Context(), contextSessionToken, token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
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
