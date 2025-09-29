package handlers

import (
	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/examples/i18n-errors/response"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var bundle = i18n.Bundle("user")

func Login(c *gin.Context) {

	var req LoginRequest
	c.ShouldBindJSON(&req)

	err := bundle.Err("User %s not exist", req.Username)

	response.Fail(c, err)
}
