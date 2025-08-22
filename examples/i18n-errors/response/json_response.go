package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 简化实现。。
func Fail(c *gin.Context, err error) {

	ctx := c.Request.Context()

	httpStatus := http.StatusInternalServerError

	var httpStatuser interface{ HttpStatus() int }
	if ok := errors.As(err, &httpStatuser); ok {
		httpStatus = httpStatuser.HttpStatus()
	}

	code := 1
	var coder interface{ Code() int }
	if ok := errors.As(err, &coder); ok {
		code = coder.Code()
	}

	var msg string
	if translatable, ok := err.(interface {
		T(ctx context.Context) string
	}); ok {
		msg = translatable.T(ctx)
	} else {
		msg = err.Error()
	}

	c.JSON(httpStatus, JsonResponse{
		Code:    code,
		Message: msg,
	})
}
