package middleware

import (
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
	"time"
)

func Log() router.Handler {
	return func(ctx *context.Context) {
		t := time.Now()
		logkit.Info("request", "url", ctx.Request.URL.String())

		ctx.Proxy.Next()

		latency := time.Since(t).Milliseconds()
		status := ctx.Proxy.Writer.Status()
		logkit.Info("response", "url", ctx.Request.URL.String(), "latency", latency, "status", status)
	}
}
