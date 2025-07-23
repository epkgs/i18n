package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var userI18n = i18n.NewCatalog("user")

var (
	ErrUserNotExit = defineErr[struct{ Name string }](userI18n, 1, "User {{.Name}} not exist", http.StatusNotFound)
)

func init() {
	userI18n.LoadTranslations()
}
