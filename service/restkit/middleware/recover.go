package middleware

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
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
					logkit.Error(e.Error())
				} else {
					msg = cast.ToString(err)
					logkit.ErrorOrigin(exception.New(msg, 3).Error())
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
