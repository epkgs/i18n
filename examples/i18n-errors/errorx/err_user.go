package errorx

import (
	"net/http"

	"github.com/epkgs/i18n/examples/i18n-errors/locales"
)

var (
	// ErrUserNotExit = Definef[string](locales.User, "User %s not exist", 1, http.StatusNotFound)
	ErrUserNotExit = Definef[struct{ Name string }](locales.User, "User {{.Name}} not exist", 1, http.StatusNotFound)
)
