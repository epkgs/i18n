package response

import (
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

	httpStatus := 500

	if e, ok := err.(*errors.Error); ok {
		httpStatus = e.HttpStatus()
		err = e.WithContext(c.Request.Context())
	}

	c.JSON(httpStatus, JsonResponse{
		Code:    1, // 非 0
		Message: err.Error(),
	})
}
