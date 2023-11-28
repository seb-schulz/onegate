package jwt

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type claims struct {
	jwt.RegisteredClaims
}

var ErrJwtInvalidSubject = errors.New("must be anon or an int")

func (m claims) Validate() error {
	sub, err := m.GetSubject()
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare([]byte(sub), []byte(AnonymousUser.Subject())) == 1 {
		return nil
	}

	uID := 0
	if _, err := fmt.Sscan(sub, "%x", &uID); err != nil {
		return ErrJwtInvalidSubject
	}
	if uID > 0 {
		return nil
	}

	return ErrJwtInvalidSubject
}

func getSecret() ([]byte, error) {
	return base64.StdEncoding.DecodeString(viper.GetString("jwt.secret"))
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(viper.GetString("jwt.header"))
		if tokenString == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access forbidden"))
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
			return getSecret()
		}, jwt.WithValidMethods(viper.GetStringSlice("jwt.valid_methods")), jwt.WithExpirationRequired(), jwt.WithLeeway(30*time.Second))

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access forbidden"))
			return
		}

		log.Println(token.Claims.GetSubject())

		next.ServeHTTP(w, r)
	})
}

func GenerateJwtToken(user UserJwtConverter) (string, error) {
	characters := "ABCDEFGHIJKLMOPQRSTUVWXYZabcdefghijklmopqrstuvwxyz0123456789"
	id_runes := make([]byte, 4)
	for i := range id_runes {
		id_runes[i] = characters[rand.Intn(len(characters))]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.expires_in"))),
		ID:        string(id_runes),
		Subject:   user.Subject(),
	})

	secret, err := getSecret()
	if err != nil {
		return "", err
	}

	return token.SignedString(secret)
}
