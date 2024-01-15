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
)

type mockEntity struct {
	data int
}

func (me *mockEntity) IDStr() string {
	return fmt.Sprint(me.data)
}

func TestStorageContext(t *testing.T) {
	sm := StorageManager[*mockEntity]{
		entityType: "mock",
	}

	for _, tc := range []struct {
		expect *mockEntity
		ctx    context.Context
	}{
		{nil, context.Background()},
		{&mockEntity{data: 1}, sm.toContext(context.Background(), &mockEntity{data: 1})},
	} {
		if got := sm.FromContext(tc.ctx); !reflect.DeepEqual(got, tc.expect) {
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

	sm := StorageManager[*mockEntity]{
		entityType: "mock",
		fetch: func(ctx context.Context) (*mockEntity, error) {
			mock, ok := fakeStorage[FromContext(ctx).UUID]
			if !ok {
				return nil, errNotFound
			}
			return mock, nil
		},
	}

	for _, tc := range scenarios {
		handler := sm.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mock := sm.FromContext(r.Context())
			if !reflect.DeepEqual(mock, tc.expectMock) {
				t.Errorf("Got mock %#v instead of %#v", mock, tc.expectMock)
			}
			fmt.Fprintln(w, "Ok")
		}))

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
