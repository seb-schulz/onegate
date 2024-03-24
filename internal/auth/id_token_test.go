package auth

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const privTestKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIEt2Ln1ml2Vq4eldfoM4QtUI8FKNr13ryEg3Hz/glRJQoAoGCCqGSM49
AwEHoUQDQgAEaQnBLC1R/I6Af5uKvwKTeAomIzHTtLZYlVYVt4U/CHsqUTtKgymq
dNmdKo8QblzagPeYu07NONnRmN5VfU3LMA==
-----END EC PRIVATE KEY-----`

const pubTestKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEaQnBLC1R/I6Af5uKvwKTeAomIzHT
tLZYlVYVt4U/CHsqUTtKgymqdNmdKo8QblzagPeYu07NONnRmN5VfU3LMA==
-----END PUBLIC KEY-----`

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
		{privKey, "https://example.com", time.Second, 1, cliendID},
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
