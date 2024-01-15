package middleware

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
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
		return config.Config.UrlLogin.Key, nil
	}, jwt.WithValidMethods(config.Config.UrlLogin.ValidMethods), jwt.WithExpirationRequired(), jwt.WithLeeway(30*time.Second))
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer http.Redirect(w, r, config.Config.BaseUrl.String(), http.StatusSeeOther)

		logger := httplog.LogEntry(r.Context())
		signedToken := chi.URLParam(r, "token")
		token, err := parseToken(signedToken)
		if err != nil {
			logger.Warn(fmt.Sprintf("cannot parse signed token: %v", err))
			return
		}

		uID, err := token.Claims.(*loginClaims).UserID()
		if err != nil {
			logger.Warn(fmt.Sprintf("cannot get user ID: %v", err))
			return
		}

		if err := model.LoginUser(r.Context(), model.LoginOpt{UserID: &uID}); err != nil {
			logger.Warn(fmt.Sprintf("cannot login usere: %v", err))
			return
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

	sigendToken, err := token.SignedString(config.Config.UrlLogin.Key)
	if err != nil {
		return nil, err

	}

	return config.Config.BaseUrl.JoinPath("login", sigendToken), nil
}
