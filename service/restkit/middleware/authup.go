package middleware

import (
	"github.com/gin-gonic/gin"
	"mizuki/project/core-kit/service/logkit"
)

/**
用户名密码校验
 */
func AuthUsernameAndPwd() gin.HandlerFunc  {
	return func(context *gin.Context) {
		logkit.Info("middleware: user: "+Session.GetString(context.Request.Context(),"me"))
		context.Next()
	}
}
