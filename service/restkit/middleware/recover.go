package middleware

import (
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
	"github.com/spf13/cast"
)

// Recover 错误处理。
func Recover() router.Handler {
	return func(ctx *context.Context) {
		defer func() {
			if err := recover(); err != nil {
				var msg string
				if e, ok := err.(exception.Exception); ok {
					msg = e.Msg
					// 带代码位置信息
					logkit.ErrorException(e)
				} else {
					msg = cast.ToString(err)
					logkit.ErrorException(exception.New(msg, 3))
				}
				if ctx.Proxy.IsAborted() {
					return
				}
				ctx.JsonError(msg)
				//ctx.Proxy.Abort()
			}
		}()
		ctx.Proxy.Next()
	}
}
