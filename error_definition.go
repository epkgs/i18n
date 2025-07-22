package i18n

import (
	"context"
)

func (i18n *I18n) DefineError(format string, wrappers ...ErrorWrapper) *errorDefinition {
	return newErrorDefinition(i18n.New(format), wrappers...)
}

type ErrorWrapper func(error) error

type errorDefinition struct {
	t        *Item // i18n item
	def      error
	wrappers []ErrorWrapper
}

func newErrorDefinition(t *Item, wrappers ...ErrorWrapper) *errorDefinition {
	return &errorDefinition{
		t:        t,
		def:      newError(t, nil),
		wrappers: wrappers,
	}
}

func (d *errorDefinition) Base() error {
	return d.def
}

func (d *errorDefinition) New(ctx context.Context, args ...any) error {

	var err error = newError(d.t, ctx, args...)
	for _, wrapper := range d.wrappers {
		err = wrapper(err)
	}

	return err
}

func (d *errorDefinition) With(wrappers ...ErrorWrapper) *errorDefinition {
	d.wrappers = append(d.wrappers, wrappers...)
	return d
}
