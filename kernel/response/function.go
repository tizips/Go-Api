package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/kernel/validator"
)

func ToResponseByUnauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40100,
		Message: "Unauthorized",
	})
}

func ToResponseByNotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40400,
		Message: message,
	})
}

func ToResponseByFailRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40000,
		Message: validator.Translate(err),
	})
}

func ToResponseByFailRequestMessage(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40000,
		Message: message,
	})
}

func ToResponseByFailLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40100,
		Message: "登陆失败",
	})
}

func ToResponseBySuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "Success",
	})
}

func ToResponseBySuccessData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "Success",
		Data:    data,
	})
}

func ToResponseBySuccessList(ctx *gin.Context, list []any) {
	ctx.JSON(http.StatusOK, Responses{
		Code:    20000,
		Message: "Success",
		Data:    list,
	})
}

func ToResponseBySuccessPaginate(ctx *gin.Context, data Paginate) {
	ctx.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "Success",
		Data:    data,
	})
}

func ToResponseByFail(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    60000,
		Message: message,
	})
}
