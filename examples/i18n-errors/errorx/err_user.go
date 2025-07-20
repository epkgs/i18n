package errorx

import (
	"net/http"
)

var userErrors = NewFactory("user")

var (
	ErrUserNotExit = userErrors.NewA1(1001, http.StatusNotFound, "User %s not exist")
)

func init() {
	userErrors.I18n.LoadTranslations()
}
