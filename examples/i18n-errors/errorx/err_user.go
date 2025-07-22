package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var userErrors = i18n.NewCatalog("user")

var (
	ErrUserNotExit = userErrors.DefineError("User %s not exist").WithStatus(1, http.StatusNotFound)
)

func init() {
	userErrors.LoadTranslations()
}
