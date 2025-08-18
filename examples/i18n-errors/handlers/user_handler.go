package handlers

import (
	"github.com/epkgs/i18n/examples/i18n-errors/errorx"
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

	err := errorx.ErrUserNotExit.New(struct{ Name string }{req.Username})

	response.Fail(c, err)
}
