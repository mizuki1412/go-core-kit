package middleware

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

// AuthUsernameAndPwd 用户名密码校验
func AuthUsernameAndPwd() router.Handler {
	return func(ctx *context.Context) {
		user := ctx.Session().Get("user")
		// todo 通过token获取session
		if user == nil {
			ctx.Json(context.RestRet{
				Result: context.ResultAuthErr,
				Message: class.String{
					String: "登录失效",
					Valid:  true,
				},
			})
			ctx.Proxy.StopExecution()
		}
		ctx.Proxy.Next()
	}
}
