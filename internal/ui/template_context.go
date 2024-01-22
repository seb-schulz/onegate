package ui

import (
	"context"
	"fmt"
	"net/http"
)

type (
	contextTemplateType struct{ string }
	templateData        map[string]any
)

var ctxTempl = contextTemplateType{"tmpl"}

func fromContext(ctx context.Context) *templateData {
	raw, ok := ctx.Value(ctxTempl).(*templateData)
	if !ok {
		return &templateData{}
	}
	return raw
}

func AddTemplateValue(ctx context.Context, key string, val any) {
	data := *fromContext(ctx)
	data[key] = val
}

func (td *templateData) Foobar() string {
	return fmt.Sprintf("%#v", td)
}

func InitTemplateContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxTempl, &templateData{"request": r})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
