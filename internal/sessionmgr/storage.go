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
		IDStr() string
	}

	StorageManager[T entity] struct {
		entityType string
		fetch      func(ctx context.Context) (T, error)
	}
)

func NewStorage[T entity](entityType string, fetchFn func(ctx context.Context) (T, error)) *StorageManager[T] {
	return &StorageManager[T]{
		entityType: entityType,
		fetch:      fetchFn,
	}
}

func (sm *StorageManager[T]) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		obj, err := sm.fetch(ctx)
		if err != nil {
			logger := httplog.LogEntry(ctx)
			logger.Info(fmt.Sprintf("cannot get %v: %v", sm.entityType, err))
		} else {
			ctx = sm.toContext(ctx, obj)
			httplog.LogEntrySetField(ctx, sm.entityType, slog.StringValue(obj.IDStr()))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (sm *StorageManager[T]) toContext(ctx context.Context, obj T) context.Context {
	return context.WithValue(ctx, ctxStorageKey{sm.entityType}, obj)

}

func (sm *StorageManager[T]) FromContext(ctx context.Context) T {
	raw, _ := ctx.Value(ctxStorageKey{sm.entityType}).(T)
	return raw

}
