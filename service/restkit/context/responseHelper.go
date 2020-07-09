package context

import (
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/storagekit"
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

func (ctx *Context) SetFileHeader(filename string) {
	ctx.Proxy.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Proxy.Header("Content-Type", "application/octet-stream")
	ctx.Proxy.Header("Content-Transfer-Encoding", "binary")
}
func (ctx *Context) FileRaw(data []byte, name string) {
	ctx.SetFileHeader(name)
	_, err := ctx.Proxy.Binary(data)
	if err != nil {
		logkit.Error("rest_ret_file_raw_error: " + err.Error())
	}
}

// 相对于项目目录路径的
func (ctx *Context) File(relativePath, name string) {
	err := ctx.Proxy.SendFile(storagekit.GetFullPath(relativePath), name)
	if err != nil {
		logkit.Error("rest_ret_file_error: " + err.Error())
	}
}
