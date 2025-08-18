package main

import (
	"context"
	"fmt"

	"github.com/epkgs/i18n"
)

func main() {
	userI18n := i18n.NewBundle("user", func(opts *i18n.Options) {
		opts.DefaultLang = "en"
	})

	userI18n.Load()

	str := userI18n.Sprintf("User %s not exist", "test")

	fmt.Printf("Default: %s\n", str)

	ctx := i18n.WithAcceptLanguages(context.Background(), "zh_CN")

	fmt.Printf("Translation [%s]: %s\n", "zh_CN", str.Translate(ctx))
}
