package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zliang90/kingRest/internal/restful/errors"
)

const (
	SuccessOK = iota // 0 means request and response ok
)

type Response struct {
	RequestId string      `json:"request_id"`
	Code      int64       `json:"code"`
	Data      interface{} `json:"data,omitempty"`
	Total     int         `json:"total,omitempty"`
}

func GetRequestId(ctx *gin.Context) string {
	return ctx.GetString("Request-Id")
}

func SuccessWithTotal(ctx *gin.Context, data interface{}, total int) {
	ctx.JSON(http.StatusOK, Response{
		RequestId: GetRequestId(ctx),
		Code:      SuccessOK,
		Data:      data,
		Total:     total,
	})
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		RequestId: GetRequestId(ctx),
		Code:      SuccessOK,
		Data:      data,
	})
}

func Failure(ctx *gin.Context, err interface{}) {
	if e, ok := err.(*errors.APIError); ok {
		e.RequestId = GetRequestId(ctx)
		ctx.AbortWithStatusJSON(http.StatusOK, e)
		return
	}
	err1 := errors.NewAPIError("UNKNOWN_ERROR", errors.Params{"error": err})
	err1.RequestId = GetRequestId(ctx)
	ctx.AbortWithStatusJSON(http.StatusOK, err1)
}
