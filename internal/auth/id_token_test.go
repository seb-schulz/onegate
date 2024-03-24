package auth

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type IDTokenConfig struct {
	key          *ecdsa.PrivateKey
	ValidMethods []string
}

const privTestKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIEt2Ln1ml2Vq4eldfoM4QtUI8FKNr13ryEg3Hz/glRJQoAoGCCqGSM49
AwEHoUQDQgAEaQnBLC1R/I6Af5uKvwKTeAomIzHTtLZYlVYVt4U/CHsqUTtKgymq
dNmdKo8QblzagPeYu07NONnRmN5VfU3LMA==
-----END EC PRIVATE KEY-----`

const pubTestKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEaQnBLC1R/I6Af5uKvwKTeAomIzHT
tLZYlVYVt4U/CHsqUTtKgymqdNmdKo8QblzagPeYu07NONnRmN5VfU3LMA==
-----END PUBLIC KEY-----`

type IdTokenClaims struct {
	jwt.RegisteredClaims
	// Nonce string `json:"nonce"`
}

type IDToken struct {
	Key       *ecdsa.PrivateKey
	Issuer    string
	ExpiresIn time.Duration
	UserID    uint
	ClientID  uuid.UUID
}

func (token IDToken) MarshalText() ([]byte, error) {
	jwt.MarshalSingleStringAsArray = false

	if token.Issuer == "" || token.UserID == 0 {
		return []byte{}, fmt.Errorf("missing values")
	}

	s := jwt.NewWithClaims(jwt.SigningMethodES256, &IdTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    token.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(token.ExpiresIn)),
			Subject:   fmt.Sprintf("%d", token.UserID),
			Audience:  jwt.ClaimStrings{fmt.Sprint(&token.ClientID)},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		// Nonce: "Nonce",
	})

	sigendToken, err := s.SignedString(token.Key)
	if err != nil {
		return []byte{}, err

	}
	return []byte(sigendToken), nil
}

func TestIDTokenMarshalJSON(t *testing.T) {
	privKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(privTestKey))
	if err != nil {
		t.Fatalf("cannot parse private test key: %v", err)
	}
	pubKey, err := jwt.ParseECPublicKeyFromPEM([]byte(pubTestKey))
	if err != nil {
		t.Fatalf("cannot parse private test key: %v", err)
	}
	cliendID := uuid.MustParse("86ec11a2-3bfc-446b-835d-35b563c10c4e")

	for _, tc := range []IDToken{
		IDToken{privKey, "https://example.com", time.Second, 1, cliendID},
	} {
		b, err := json.Marshal(tc)
		if err != nil {
			t.Errorf("failed to marshal token: %v", err)
		}

		parsedToken, err := jwt.ParseWithClaims(string(b[1:len(b)-1]), &IdTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		}, jwt.WithValidMethods([]string{"ES256"}), jwt.WithExpirationRequired(), jwt.WithLeeway(30*time.Second))
		if err != nil {
			t.Errorf("failed to parse token: %v", err)
		}
		t.Log(parsedToken)
		t.Log(parsedToken.Claims)

		if sub, _ := parsedToken.Claims.GetSubject(); sub != fmt.Sprint(tc.UserID) {
			t.Errorf("Expected %d but got %s", tc.UserID, sub)
		}
		if iss, _ := parsedToken.Claims.GetIssuer(); iss != tc.Issuer {
			t.Errorf("Expected %s but got %s", tc.Issuer, iss)
		}
		if aud, _ := parsedToken.Claims.GetAudience(); aud[0] != fmt.Sprint(tc.ClientID) {
			t.Errorf("Expected %s but got %s", tc.ClientID, aud)
		}

	}
}

func newIDToken(ic IDTokenConfig) (string, error) {
	jwt.MarshalSingleStringAsArray = false

	token := jwt.NewWithClaims(jwt.SigningMethodES256, &IdTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://example.com",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			Subject:   fmt.Sprintf("%d", 1),
			Audience:  jwt.ClaimStrings{"1234"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		// Nonce: "Nonce",
	},
	)

	sigendToken, err := token.SignedString(ic.key)
	if err != nil {
		return "", err

	}
	return sigendToken, nil
}

func TestNewIDToken(t *testing.T) {
	t.Skip()
	privKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(privTestKey))
	if err != nil {
		t.Fatalf("cannot parse private test key: %v", err)
	}

	token, err := newIDToken(IDTokenConfig{
		privKey,
		[]string{"HS256"},
	})
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	t.Log(token)

	t.FailNow()
}
