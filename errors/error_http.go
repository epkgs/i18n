package errors

import (
	"net/http"

	"github.com/epkgs/i18n/types"
)

// Default error code and HTTP status
const (
	CodeDefault       = 1
	HttpStatusDefault = http.StatusInternalServerError
)

// WithCode sets a custom error code on the error
//   - E must implement the Storager interface to store the code
//
// Returns the same error with the code set
func WithCode[E types.Storager](err E, code int) E {
	err.Set("code", code)
	return err
}

// Code retrieves the error code from an error
//   - If the error is nil, returns the default code
//   - If the error implements Storager and has a code set, returns that code
//   - If the error has a Code() method, returns the result of that method
//
// Otherwise, returns the default code
func Code(err error) int {
	if err == nil {
		return CodeDefault
	}

	if storage, ok := err.(types.Storager); ok {

		if code, ok := storage.Get("code", CodeDefault).(int); ok {
			return code
		}

		return CodeDefault
	}

	var coder interface{ Code() int }
	if ok := As(err, &coder); ok {
		return coder.Code()
	}

	return CodeDefault
}

// WithHttpStatus sets a custom HTTP status code on the error
// E must implement the Storager interface to store the HTTP status
// Returns the same error with the HTTP status set
func WithHttpStatus[E types.Storager](err E, httpStatus int) E {
	err.Set("http_status", httpStatus)
	return err
}

// HttpStatus retrieves the HTTP status code from an error
//   - If the error is nil, returns the default HTTP status
//   - If the error implements Storager and has an HTTP status set, returns that status
//   - If the error has an HttpStatus() method, returns the result of that method
//
// Otherwise, returns the default HTTP status
func HttpStatus(err error) int {
	if err == nil {
		return HttpStatusDefault
	}

	if storage, ok := err.(types.Storager); ok {
		if status, ok := storage.Get("http_status", HttpStatusDefault).(int); ok {
			return status
		}
		return HttpStatusDefault
	}

	var httpStatus interface{ HttpStatus() int }
	if ok := As(err, &httpStatus); ok {
		return httpStatus.HttpStatus()
	}

	return HttpStatusDefault
}
