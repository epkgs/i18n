package response

import (
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
	{
		if e, ok := err.(interface{ HttpStatus() int }); ok {
			httpStatus = e.HttpStatus()
		}
	}

	code := 1
	{
		if e, ok := err.(interface{ Code() int }); ok {
			code = e.Code()
		}
	}

	c.JSON(httpStatus, JsonResponse{
		Code:    code,
		Message: err.Error(),
	})
}
