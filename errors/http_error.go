package errors

import (
	"errors"
	"fmt"
)

var (
	CodeDefault       = 1
	HttpStatusDefault = 500
)

func NewHttpError(code, httpStatus int, format fmt.Stringer) I18nError {
	err := New(format)
	err.WithCode(code)
	err.WithHttpStatus(httpStatus)
	return err
}

func (e *Error) WithCode(code int) I18nError {
	err := e.clone()
	err.Set("code", code)
	return err
}

func (e *Error) Code() int {
	if code, exist := e.extra["code"]; exist {
		return code.(int)
	}

	return CodeDefault
}

func (e *Error) WithHttpStatus(httpStatus int) I18nError {
	err := e.clone()
	err.Set("httpStatus", httpStatus)
	return err
}

func (e *Error) HttpStatus() int {
	if httpStatus, exist := e.extra["httpStatus"]; exist {
		return httpStatus.(int)
	}

	return HttpStatusDefault
}

func Code(err error) int {
	type coder interface{ Code() int }
	if err == nil {
		return CodeDefault
	}

	var c coder
	if ok := errors.As(err, &c); ok {
		return c.Code()
	}

	return CodeDefault
}

func HttpStatus(err error) int {
	type httpStatus interface{ HttpStatus() int }
	if err == nil {
		return HttpStatusDefault
	}

	var s httpStatus
	if ok := errors.As(err, &s); ok {
		return s.HttpStatus()
	}

	return HttpStatusDefault
}
