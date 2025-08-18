package errors

import "github.com/epkgs/i18n"

type Definition[Args any] struct {
	i18n   *i18n.Bundle
	format string

	opts []DefineOption
}

type DefineOption func(e I18nError)

func Definef[Args any](i18n *i18n.Bundle, format string, opts ...DefineOption) *Definition[Args] {
	return &Definition[Args]{
		i18n:   i18n,
		format: format,
		opts:   opts,
	}
}

func (d *Definition[Args]) New(args Args) I18nError {
	err := New(d.i18n.Sprintf(d.format, args))

	for _, opt := range d.opts {
		opt(err)
	}

	return err
}
