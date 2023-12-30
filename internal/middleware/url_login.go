package middleware

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/model"
	"gorm.io/gorm"
)

type loginClaims struct {
	jwt.RegisteredClaims
}

var errJwtInvalidSubject = errors.New("must be an int greater than zero")

func (m loginClaims) Validate() error {
	uID, err := m.UserID()
	if err != nil {
		return err
	}

	if uID <= 0 {
		return errJwtInvalidSubject
	}

	return nil
}

func (m loginClaims) UserID() (uint, error) {
	sub, err := m.GetSubject()
	if err != nil {
		return 0, err
	}

	uID := uint(0)
	if _, err := fmt.Sscanf(sub, "%x", &uID); err != nil {
		return uID, errJwtInvalidSubject
	}
	return uID, nil
}

func parseToken(signedToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(signedToken, &loginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Default.UrlLogin.Key, nil
	}, jwt.WithValidMethods(config.Default.JWT.ValidMethods), jwt.WithExpirationRequired(), jwt.WithLeeway(30*time.Second))
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer http.Redirect(w, r, config.Default.BaseUrl.String(), http.StatusSeeOther)

		signedToken := chi.URLParam(r, "token")
		token, err := parseToken(signedToken)
		if err != nil {
			log.Println(err)
			return
		}

		uID, err := token.Claims.(*loginClaims).UserID()
		if err != nil {
			log.Println(err)
			return
		}
		session := SessionFromContext(r.Context())

		if err := db.Transaction(func(tx *gorm.DB) error {
			user := model.User{}
			tx.First(&user, "id = ?", uID)

			session.User = &user
			tx.Save(session)

			return nil
		}); err != nil {
			panic(err)
		}

		// w.Write([]byte(`<!DOCTYPE html>
		// <html>
		// <head><meta http-equiv="refresh" content="0; url='/'"></head>
		// <body></body>
		// </html>`))
	})
}

func GetLoginUrl(userID uint, expiresIn time.Duration) (*url.URL, error) {
	characters := "ABCDEFGHIJKLMOPQRSTUVWXYZabcdefghijklmopqrstuvwxyz0123456789"
	id_runes := make([]byte, 4)
	for i := range id_runes {
		id_runes[i] = characters[rand.Intn(len(characters))]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		ID:        string(id_runes),
		Subject:   fmt.Sprintf("%x", userID),
	})

	sigendToken, err := token.SignedString(config.Default.UrlLogin.Key)
	if err != nil {
		return nil, err

	}

	return config.Default.BaseUrl.JoinPath("login", sigendToken), nil
}
