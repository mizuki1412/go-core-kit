package middleware

import (
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"time"
)

func Log() router.Handler {
	return func(ctx *context.Context) {
		t := time.Now()
		// 当前上传的cookies中的session，不一定等于response中的
		sessionId := ctx.SessionID()
		// todo params
		//params := make(map[interface{}]interface{})
		//_ = c.Copy().ShouldBind(&params)
		logkit.Info("request",
			logkit.Param{
				Key: "session",
				Val: sessionId,
			}, logkit.Param{
				Key: "url",
				Val: ctx.Request.URL.String(),
			})

		ctx.Proxy.Next()

		latency := time.Since(t).Milliseconds()
		status := ctx.Proxy.Writer.Status()
		logkit.Info("response",
			logkit.Param{
				Key: "url",
				Val: ctx.Request.URL.String(),
			}, logkit.Param{
				Key: "latency",
				Val: latency,
			}, logkit.Param{
				Key: "status",
				Val: status,
			})
	}
}
