package middleware

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"github.com/spf13/cast"
)

/** 错误处理，以及db事务处理。 */
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
					logkit.Error(exception.New(msg, 3))
				}
				// transaction
				if ctx.DBTxExist() {
					ctx.DBTx().Rollback()
				}
				if ctx.Proxy.IsStopped() {
					return
				}
				ctx.JsonError(msg)
				// 打印错误信息
				logkit.Info("response-error",
					logkit.Param{
						Key: "url",
						Val: ctx.Request.URL.String(),
					}, logkit.Param{
						Key: "msg",
						Val: msg,
					})
				ctx.Proxy.StopExecution()
			}
		}()
		ctx.Proxy.Next()
		// transaction
		if ctx.DBTxExist() {
			ctx.DBTx().Commit()
		}
	}
}
