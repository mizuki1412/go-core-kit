package middleware

import (
	"github.com/spf13/cast"
	"mizuki/project/core-kit/class/exception"
	"mizuki/project/core-kit/library/stringkit"
	"mizuki/project/core-kit/service/logkit"
	"mizuki/project/core-kit/service/restkit/context"
	"mizuki/project/core-kit/service/restkit/router"
)

/**
错误处理，以及db事务处理。
*/
func Recover() router.Handler {
	return func(ctx *context.Context) {
		defer func() {
			if err := recover(); err != nil {
				var msg string
				if e, ok := err.(exception.Exception); ok {
					msg = e.Error()
					logkit.Error(e.Error(), logkit.Param{
						Key: "position",
						Val: stringkit.Concat(e.File, ":", cast.ToString(e.Line)),
					})
				} else {
					msg = cast.ToString(err)
					logkit.Error(msg)
				}
				// transaction
				if ctx.DBTxExist() {
					ctx.DBTx().Rollback()
				}
				if ctx.Proxy.IsStopped() {
					return
				}
				ctx.JsonError(msg)
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
