package main

import (
	"context"
	"fmt"

	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/examples/simple/locales"
)

//go:generate go run ../../i18ntool extract

func main() {

	ctx := context.Background()

	user := "test"

	userNotExist := locales.User.Str("User %s not exist", user)

	fmt.Printf("Default: %s\n", userNotExist)

	ctx = i18n.WithAcceptLanguages(ctx, "zh_CN")

	fmt.Printf("Translation [zh_CN]: %s\n", userNotExist.T(ctx))
}
