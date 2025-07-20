package errorx

import (
	"net/http"

	"github.com/epkgs/i18n/errors"
)

var userErrors = errors.NewBuilder("user")

var (
	ErrUserNotExit = userErrors.New(1, "User %s not exist", http.StatusNotFound)
)

func init() {
	userErrors.I18n.LoadTranslations()
}
