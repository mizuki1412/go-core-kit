package restkit

import (
	"mizuki/project/core-kit/service/logkit"
	"mizuki/project/core-kit/service/restkit/context"
	"net/http"
)

type RestRet struct {
	Result  int         `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const RestRetResultErr = 0
const RestRetResultSuccess = 1
const RestRetResultAuthErr = 2

// http返回json数据
func RetJson(context *context.Context, ret RestRet) {
	var code int
	switch ret.Result {
	case RestRetResultSuccess:
		code = http.StatusOK
	case RestRetResultAuthErr:
		code = http.StatusUnauthorized
	default:
		code = http.StatusBadRequest
	}
	context.Proxy.StatusCode(code)
	_, err := context.Proxy.JSON(ret)
	if err != nil {
		logkit.Error("rest_ret_json_error: " + err.Error())
	}
}
