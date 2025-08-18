package response

import (
	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/examples/i18n-errors/errorx"
	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 简化实现。。
func Fail(c *gin.Context, err error) {

	httpStatus := errorx.HttpStatus(err)

	var msg string
	if e, ok := err.(i18n.Translable); ok {
		msg = e.Translate(c.Request.Context())
	} else {
		msg = err.Error()
	}

	c.JSON(httpStatus, JsonResponse{
		Code:    errorx.Code(err),
		Message: msg,
	})
}
