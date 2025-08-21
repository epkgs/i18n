package i18n

import (
	"context"
	"errors"
)

type ErrorDefinitionf[Args any] struct {
	base *Definitionf[Args]
}

func DefineErrorf[Args any](i18n *Bundle, other string, one ...string) *ErrorDefinitionf[Args] {
	return &ErrorDefinitionf[Args]{
		base: Definef[Args](i18n, other, one...),
	}
}

func (e *ErrorDefinitionf[Args]) T(ctx context.Context, args Args, count ...any) error {
	return errors.New(e.base.T(ctx, args, count...))
}

type ErrorDefinition struct {
	base *Definition
}

func DefineError(i18n *Bundle, txt string) *ErrorDefinition {
	return &ErrorDefinition{
		base: Define(i18n, txt),
	}
}

func (e *ErrorDefinition) T(ctx context.Context) error {
	return errors.New(e.base.T(ctx))
}

func (e *ErrorDefinition) String() string {
	return e.base.String()
}
