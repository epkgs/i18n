package main

import (
	"context"
	"fmt"

	"github.com/epkgs/i18n"
)

func main() {
	userI18n := i18n.New("user", func(opts *i18n.Options) {
		opts.DefaultLang = "en"
	})

	UserNotExist := userI18n.NewItem("User %s not exist")

	i18n.LoadTranslations(userI18n)

	fmt.Printf("Default: %s\n", UserNotExist.T(context.Background(), "test"))

	ctx := i18n.WithAcceptLanguages(context.Background(), "zh_CN")
	fmt.Printf("Translation [%s]: %s\n", "zh_CN", UserNotExist.T(ctx, "test"))
}
