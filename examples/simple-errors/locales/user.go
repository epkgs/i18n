package locales

import "github.com/epkgs/i18n"

var userI18n = i18n.New("user")

var (
	UserNotExist = userI18n.New("User %s not exist")
)

func init() {
	i18n.LoadTranslations(userI18n)
}
