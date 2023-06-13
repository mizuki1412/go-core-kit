package middleware

import (
	"github.com/mizuki1412/go-core-kit/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"github.com/spf13/cast"
	"time"
)

// CreateSession 在登录等需要开启session的时候使用
func CreateSession() router.Handler {
	return func(ctx *context.Context) {
		token := context.GetTokenFromReq(ctx)
		if token == "" {
			// 开启session id
			token = cryptokit.ID() + "-" + cast.ToString(time.Now().UnixMilli())
		}
		if ctx.Get("_token") == nil {
			ctx.Set("_token", token)
		}
		ctx.Proxy.Next()
	}
}
