package handlers

import (
	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/examples/i18n-errors/locales"
	"github.com/epkgs/i18n/examples/i18n-errors/response"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {

	var req LoginRequest
	c.ShouldBindJSON(&req)

	def := i18n.DefineErrorf[string](locales.User, "User %s not exist")

	err := def.T(c.Request.Context(), req.Username)

	response.Fail(c, err)
}
