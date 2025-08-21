package main

import (
	"context"
	"fmt"

	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/examples/simple/locales"
)

//go:generate go run ../../geni18n extract

func main() {

	ctx := context.Background()

	userNotExist := i18n.Definef[string](locales.User, "User %s not exist")

	fmt.Printf("Default: %s\n", userNotExist.T(ctx, "test"))

	ctx = i18n.WithAcceptLanguages(ctx, "zh_CN")

	fmt.Printf("Translation [zh_CN]: %s\n", userNotExist.T(ctx, "test"))
}
