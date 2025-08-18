package response

import (
	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/errors"
	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 简化实现。。
func Fail(c *gin.Context, err error) {

	httpStatus := errors.HttpStatus(err)

	var msg string
	if e, ok := err.(i18n.Translable); ok {
		msg = e.Translate(c.Request.Context())
	} else {
		msg = err.Error()
	}

	c.JSON(httpStatus, JsonResponse{
		Code:    errors.Code(err),
		Message: msg,
	})
}
