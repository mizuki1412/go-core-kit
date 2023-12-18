package middleware

import (
	"github.com/mizuki1412/go-core-kit/service/cachekit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"time"
)

// AuthJWT 用户名密码校验
func AuthJWT() router.Handler {
	return func(ctx *context.Context) {
		jwt := ctx.GetJwt()
		// 获取 jwt
		if !jwt.IsValid() || jwt.ExpiresAt.Before(time.Now()) || cachekit.Get("token:"+jwt.Token()) == nil {
			ctx.Json(context.RestRet{
				Result:  context.ResultAuthErr,
				Message: "登录失效",
			})
			ctx.Proxy.Abort()
		} else {
			ctx.Proxy.Next()
		}
	}
}
