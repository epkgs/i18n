package errorx

import (
	"net/http"

	"github.com/epkgs/i18n/errors"
)

var userErrors = errors.New("user")

var (
	ErrUserNotExit = userErrors.New("User %s not exist").WithHttpStatus(http.StatusNotFound)
)

func init() {
	userErrors.I18n.LoadTranslations()
}
