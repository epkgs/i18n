package i18n

import (
	"context"
)

type (
	acceptLanguagesCtx struct{}
)

func WithAcceptLanguages(ctx context.Context, acceptLanguages ...string) context.Context {
	return context.WithValue(ctx, acceptLanguagesCtx{}, acceptLanguages)
}

func GetAcceptLanguages(ctx context.Context) []string {
	v := ctx.Value(acceptLanguagesCtx{})
	if v != nil {
		return v.([]string)
	}
	return nil
}
