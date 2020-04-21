package middleware

import (
	"mizuki/project/core-kit/service/logkit"
	"mizuki/project/core-kit/service/restkit/context"
	"time"
)

func Log() context.Handler {
	return func(ctx *context.Context) {
		t := time.Now()
		// 当前上传的cookies中的session，不一定等于response中的
		sessionId := ctx.Session().ID()
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
		status := ctx.Response.StatusCode()
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