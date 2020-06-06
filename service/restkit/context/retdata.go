package context

import (
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"net/http"
)

type RestRet struct {
	Result  int         `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const ResultErr = 0
const ResultSuccess = 1
const ResultAuthErr = 2

// http返回json数据
func (ctx *Context) Json(ret RestRet) {
	var code int
	switch ret.Result {
	case ResultSuccess:
		code = http.StatusOK
	case ResultAuthErr:
		code = http.StatusUnauthorized
	default:
		code = http.StatusBadRequest
	}
	ctx.Proxy.StatusCode(code)
	_, err := ctx.Proxy.JSON(ret)
	if err != nil {
		logkit.Error("rest_ret_json_error: " + err.Error())
	}
}

func (ctx *Context) JsonSuccess(data interface{}) {
	// todo 更新session的expire 会不会太频繁
	ctx.UpdateSessionExpire()
	ctx.Json(RestRet{
		Result: ResultSuccess,
		Data:   data,
	})
}
func (ctx *Context) JsonError(msg string) {
	ctx.Json(RestRet{
		Result:  ResultErr,
		Message: msg,
	})
}
