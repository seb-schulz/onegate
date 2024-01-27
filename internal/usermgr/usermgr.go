package usermgr

import (
	"context"
	"net/http"

	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/sessionmgr"
)

var (
	defaultMgr *sessionmgr.StorageManager[*model.User]
)

func init() {
	defaultMgr = sessionmgr.NewStorage("user", model.FirstUser)
}

func Middleware(next http.Handler) http.Handler {
	return defaultMgr.Handler(next)
}

func FromContext(ctx context.Context) *model.User {
	return defaultMgr.FromContext(ctx)
}
