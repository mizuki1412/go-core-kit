package middleware

import (
	"time"

	"github.com/mizuki1412/go-core-kit/v2/service/cachekit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
)

// AuthJWT 用户名密码校验
func AuthJWT() router.Handler {
	return func(ctx *context.Context) {
		jwt := ctx.GetJwt()
		// 获取 jwt
		if !jwt.IsValid() || jwt.ExpiresAt.Before(time.Now()) || cachekit.Get("token:"+jwt.Token()) == "" {
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
