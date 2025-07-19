package errorx

import (
	"net/http"
)

var userFactory = NewFactory("user")

var (
	ErrUserNotExit = userFactory.NewA1(1001, http.StatusNotFound, "User %s not exist")
)

func init() {
	userFactory.I18n.LoadLocales()
}
