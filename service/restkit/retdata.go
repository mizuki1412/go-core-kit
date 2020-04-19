package restkit

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RestRet struct {
	Result int		`json:"result"`
	Message string	`json:"message"`
	Data interface{} `json:"data"`
}

const RestRetResultErr = 0
const RestRetResultSuccess = 1
const RestRetResultAuthErr = 2

// http返回json数据
func RetJson(context *gin.Context, ret RestRet) {
	var code int
	switch ret.Result {
	case RestRetResultSuccess:
		code = http.StatusOK
	case RestRetResultAuthErr:
		code = http.StatusUnauthorized
	default:
		code = http.StatusBadRequest
	}
	context.JSON(code, ret)
}