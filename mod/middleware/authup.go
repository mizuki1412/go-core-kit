package middleware

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

// AuthUsernameAndPwd 用户名密码校验
func AuthUsernameAndPwd() router.Handler {
	return func(ctx *context.Context) {
		// 注意另外的token情况是在 create_session
		token := context.GetTokenFromReq(ctx)
		if token != "" && ctx.Get("_token") == nil {
			ctx.Set("_token", token)
		}
		user := ctx.SessionGetUserOrigin()
		if user == nil {
			ctx.Json(context.RestRet{
				Result: context.ResultAuthErr,
				Message: class.String{
					String: "登录失效",
					Valid:  true,
				},
			})
			ctx.Proxy.Abort()
		} else {
			ctx.SessionRenew()
			ctx.Proxy.Next()
		}
	}
}
