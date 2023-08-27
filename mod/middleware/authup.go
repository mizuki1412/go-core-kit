package middleware

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

// AuthUsernameAndPwd 用户名密码校验
func AuthUsernameAndPwd() router.Handler {
	return func(ctx *context.Context) {
		// 获取 jwt
		if !ctx.GetJwt().Valid {
			ctx.Json(context.RestRet{
				Result:  context.ResultAuthErr,
				Message: class.NewString("登录失效"),
			})
			ctx.Proxy.Abort()
		} else {
			ctx.Proxy.Next()
		}
	}
}
