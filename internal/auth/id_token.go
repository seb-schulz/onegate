package auth

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

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

	if token.Issuer == "" {
		return []byte{}, fmt.Errorf("missing issuer")
	}

	if token.UserID == 0 {
		return []byte{}, fmt.Errorf("missing user ID")
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
