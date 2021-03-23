package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewWebError(ctx *gin.Context, code int, error string) {
	ctx.JSON(code, WebResponse{
		Code:  code,
		Error: error,
	})
}

func NewWebResponse(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, WebResponse{
		Code:    code,
		Message: message,
	})
}

func NewBadRequestError(ctx *gin.Context, error string) {
	NewWebError(ctx, http.StatusBadRequest, error)
}

func NewNotFoundError(ctx *gin.Context, error string) {
	NewWebError(ctx, http.StatusNotFound, error)
}

func NewOkResponse(ctx *gin.Context, message string) {
	NewWebResponse(ctx, http.StatusOK, message)
}
