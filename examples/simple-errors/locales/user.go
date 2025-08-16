package locales

import "github.com/epkgs/i18n"

var userI18n = i18n.NewBundle("user")

var (
	UserNotExist = userI18n.Define("User %s not exist")
)

func init() {
	i18n.Load(userI18n)
}
