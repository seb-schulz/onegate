package middleware

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/seb-schulz/onegate/internal/model"
)

type (
	loginClaims struct {
		jwt.RegisteredClaims
	}

	tokenBasedLoginService struct {
		key          []byte
		validMethods []string
		baseUrl      url.URL
		loginUser    func(context.Context, model.LoginOpt) error
		targetUrl    url.URL
	}

	LoginConfig struct {
		Key          []byte
		ValidMethods []string
		BaseUrl      url.URL
	}
)

var (
	errJwtInvalidSubject = errors.New("must be an int greater than zero")
	defaultTargetUrl     = url.URL{Path: "/"}
)

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

func NewLoginRoute(lc LoginConfig) http.Handler {
	route := chi.NewRouter()

	tokenSrv := &tokenBasedLoginService{lc.Key, lc.ValidMethods, lc.BaseUrl, model.LoginUser, defaultTargetUrl}
	route.Get("/{token}", tokenSrv.Handler)

	return route
}

func TokenBasedLoginUrl(lc LoginConfig, userID uint, expiresIn time.Duration) (*url.URL, error) {
	tokenSrv := &tokenBasedLoginService{lc.Key, lc.ValidMethods, lc.BaseUrl, model.LoginUser, defaultTargetUrl}
	return tokenSrv.getLoginUrl(userID, expiresIn)
}

func (ls *tokenBasedLoginService) parseToken(signedToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(signedToken, &loginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return ls.key, nil
	}, jwt.WithValidMethods(ls.validMethods), jwt.WithExpirationRequired(), jwt.WithLeeway(30*time.Second))
}

func (ls *tokenBasedLoginService) getLoginUrl(userID uint, expiresIn time.Duration) (*url.URL, error) {
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

	sigendToken, err := token.SignedString(ls.key)
	if err != nil {
		return nil, err

	}

	return ls.baseUrl.JoinPath(sigendToken), nil
}

func (ls *tokenBasedLoginService) Handler(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, fmt.Sprint(ls.targetUrl), http.StatusSeeOther)

	logger := httplog.LogEntry(r.Context())
	signedToken := chi.URLParam(r, "token")
	token, err := ls.parseToken(signedToken)
	if err != nil {
		logger.Warn(fmt.Sprintf("cannot parse signed token: %v", err))
		return
	}

	uID, err := token.Claims.(*loginClaims).UserID()
	if err != nil {
		logger.Warn(fmt.Sprintf("cannot get user ID: %v", err))
		return
	}

	if err := ls.loginUser(r.Context(), model.LoginOpt{UserID: &uID}); err != nil {
		logger.Warn(fmt.Sprintf("cannot login usere: %v", err))
		return
	}

	// w.Write([]byte(`<!DOCTYPE html>
	// <html>
	// <head><meta http-equiv="refresh" content="0; url='/'"></head>
	// <body></body>
	// </html>`))
}
