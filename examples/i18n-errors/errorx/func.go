package errorx

import (
	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/errors"
)

func Define[Args any](i18n *i18n.Bundle, format string, code, httpStatus int) *errors.Definition[Args] {
	return errors.Define[Args](i18n, format, func(e errors.I18nError) {
		e.WithCode(code)
		e.WithHttpStatus(httpStatus)
	})
}

func New(i18n *i18n.Bundle, format string, code, httpStatus int) errors.I18nError {
	err := errors.New(i18n.Sprintf(format))
	err.WithCode(code)
	err.WithHttpStatus(httpStatus)
	return err
}
