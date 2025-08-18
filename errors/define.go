package errors

import "github.com/epkgs/i18n"

type definition[Args any] struct {
	i18n   *i18n.Bundle
	format string

	opts []DefineOption
}

type DefineOption func(e I18nError) I18nError

func Define[Args any](i18n *i18n.Bundle, format string, opts ...DefineOption) Definition[Args] {
	return &definition[Args]{
		i18n:   i18n,
		format: format,
		opts:   opts,
	}
}

func (d *definition[Args]) New(args Args) I18nError {
	err := New(d.i18n.Sprintf(d.format, args))

	for _, opt := range d.opts {
		opt(err)
	}

	return err
}
