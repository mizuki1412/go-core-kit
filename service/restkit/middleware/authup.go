package middleware

import (
	"mizuki/project/core-kit/service/logkit"
	"mizuki/project/core-kit/service/restkit/context"
)

/**
用户名密码校验
*/
func AuthUsernameAndPwd() context.Handler {
	return func(context *context.Context) {
		logkit.Info("middleware: user: " + context.Session().GetString("me"))
		context.Proxy.Next()
	}
}
