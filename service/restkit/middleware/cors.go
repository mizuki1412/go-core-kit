package middleware

import (
	"mizuki/project/core-kit/service/restkit/context"
)

func Cors() context.Handler {
	return func(c *context.Context) {
		//method := c.Request.Method
		c.Proxy.Header("Access-Control-Allow-Origin", "*")
		c.Proxy.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Proxy.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Proxy.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Proxy.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		//if method == "OPTIONS" {
		//	c.Proxy.StatusCode(http.StatusNoContent)
		//	c.Proxy.StopExecution()
		//}
		c.Proxy.Next()
	}
}
