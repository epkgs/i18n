package response

import (
	"errors"

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

		var e interface{ HttpStatus() int }
		if errors.As(err, &e) {
			httpStatus = e.HttpStatus()
		}
	}

	code := 1
	{
		var e interface{ Code() int }
		if errors.As(err, &e) {
			code = e.Code()
		}
	}

	c.JSON(httpStatus, JsonResponse{
		Code:    code,
		Message: err.Error(),
	})
}
