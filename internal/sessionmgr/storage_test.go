package sessionmgr

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	middleware_logger "github.com/seb-schulz/onegate/internal/middleware"
)

type mockEntity struct {
	data int
}

func (me *mockEntity) String() string {
	return fmt.Sprint(me.data)
}

func TestStorageContext(t *testing.T) {
	sm := storageManager[*mockEntity]{
		entityType: "mock",
	}

	for _, tc := range []struct {
		expect *mockEntity
		ctx    context.Context
	}{
		{nil, context.Background()},
		{&mockEntity{data: 1}, sm.toContext(context.Background(), &mockEntity{data: 1})},
	} {
		if got := sm.fromContext(tc.ctx); !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("got %#v instead of %#v", got, tc.expect)
		}
	}
}

func TestStorageFetcher(t *testing.T) {
	type testCase struct {
		expectMock *mockEntity
		ctx        context.Context
	}

	var (
		fakeStorage = make(map[uuid.UUID]*mockEntity)
		scenarios   []testCase
		errNotFound = fmt.Errorf("not found")
	)

	seed := rand.Int()
	gen := rand.New(rand.NewSource(int64(seed)))

	newRandToken := func() *Token {
		id, _ := uuid.NewRandomFromReader(gen)
		return &Token{UUID: id}

	}

	token := newRandToken()
	mockID := gen.Int()
	mock := mockEntity{data: mockID}

	scenarios = append(scenarios, testCase{
		expectMock: &mock,
		ctx:        context.WithValue(context.Background(), contextToken, token),
	})
	fakeStorage[token.UUID] = &mock

	token = newRandToken()
	scenarios = append(scenarios, testCase{
		expectMock: nil,
		ctx:        context.WithValue(context.Background(), contextToken, token),
	})

	sm := storageManager[*mockEntity]{
		entityType: "mock",
		fetch: func(t *Token) (*mockEntity, error) {
			mock, ok := fakeStorage[t.UUID]
			if !ok {
				return nil, errNotFound
			}
			return mock, nil
		},
	}

	for _, tc := range scenarios {
		handler := sm.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mock := sm.fromContext(r.Context())
			if !reflect.DeepEqual(mock, tc.expectMock) {
				t.Errorf("Got mock %#v instead of %#v", mock, tc.expectMock)
			}
			fmt.Fprintln(w, "Ok")
		}))
		handler = middleware_logger.Logger(handler)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, newCustomRequest(func(r *http.Request) {
			*r = *r.WithContext(tc.ctx)
		}))

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.FailNow()
		}
	}
}
