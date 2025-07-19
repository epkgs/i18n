package errorx

import "net/http"

var (
	ErrNotFound = New(http.StatusNotFound, http.StatusText(http.StatusNotFound))
)
