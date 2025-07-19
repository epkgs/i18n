package handlers

import (
	"github.com/epkgs/i18n/examples/simple-errors/errorx"
	"github.com/epkgs/i18n/examples/simple-errors/locales"
	"github.com/epkgs/i18n/examples/simple-errors/response"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {

	var req LoginRequest
	c.ShouldBindJSON(&req)

	ctx := c.Request.Context()

	err := errorx.ErrNotFound.WithMsg(locales.UserNotExist.T(ctx, req.Username))

	response.Fail(c, err)
}
