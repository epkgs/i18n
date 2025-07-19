package main

import (
	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/examples/simple-errors/handlers"
	_ "github.com/epkgs/i18n/examples/simple-errors/locales"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func main() {
	r := gin.Default()

	r.Use(i18n.GinMiddleware(language.AmericanEnglish.String()))

	r.POST("/api/v1/user/login", handlers.Login)

	r.Run(":8080")

}
