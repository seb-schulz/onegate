package sessionmgr

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

type (
	ctxStorageKey struct{ string }

	entity interface {
		fmt.Stringer
	}

	storageManager[T entity] struct {
		entityType string
		fetch      func(*Token) (T, error)
	}
)

func (sm *storageManager[T]) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		obj, err := sm.fetch(FromContext(ctx))
		if err != nil {
			logger := httplog.LogEntry(ctx)
			logger.Info(fmt.Sprintf("cannot get %v: %v", sm.entityType, err))
		}

		ctx = sm.toContext(ctx, obj)
		httplog.LogEntrySetField(ctx, sm.entityType, slog.StringValue(fmt.Sprint(obj)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (sm *storageManager[T]) toContext(ctx context.Context, obj T) context.Context {
	return context.WithValue(ctx, ctxStorageKey{sm.entityType}, obj)

}

func (sm *storageManager[T]) fromContext(ctx context.Context) T {
	raw, _ := ctx.Value(ctxStorageKey{sm.entityType}).(T)
	return raw

}
