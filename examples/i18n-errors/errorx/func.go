package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/errors"
)

const (
	CodeFail    = 1
	CodeSuccess = 0
)

func Define[Args any](i18n *i18n.Bundle, format string, code, httpStatus int) errors.Definition[Args] {
	return errors.Define[Args](i18n, format, func(e errors.I18nError) errors.I18nError {
		e.Set("code", code)
		e.Set("httpStatus", httpStatus)
		return e
	})
}

func New(i18n *i18n.Bundle, format string, code, httpStatus int) errors.I18nError {
	err := errors.New(i18n.Sprintf(format))
	err.Set("code", code)
	err.Set("httpStatus", httpStatus)
	return err
}

func Code(err error) int {
	if err == nil {
		return CodeFail
	}

	var getter interface{ Get(key string) (any, bool) }
	if ok := errors.As(err, &getter); ok {
		if code, ok := getter.Get("code"); ok {
			if c, ok := code.(int); ok {
				return c
			}
		}
	}

	var coder interface{ Code() int }
	if ok := errors.As(err, &coder); ok {
		return coder.Code()
	}

	return CodeFail
}

func HttpStatus(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	var getter interface{ Get(key string) (any, bool) }
	if ok := errors.As(err, &getter); ok {
		if httpStatus, ok := getter.Get("httpStatus"); ok {
			if c, ok := httpStatus.(int); ok {
				return c
			}
		}
	}

	var httpStatuser interface{ HttpStatus() int }
	if ok := errors.As(err, &httpStatuser); ok {
		return httpStatuser.HttpStatus()
	}

	return http.StatusInternalServerError
}
